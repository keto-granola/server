package middleware

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/keto-granola/server/internal/config"
)

func Log(next echo.HandlerFunc) echo.HandlerFunc {
	return func(e echo.Context) error {
		start := time.Now()
		req := e.Request()

		var bodyBytes []byte
		if req.Body != nil {
			// no need to handle error, if it fails bodyBytes will be empty and then won't be logged
			// request shouldn't be failed over a logging issue
			bodyBytes, _ = io.ReadAll(req.Body)

			req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		allParams := getRequestParams(e)

		err := next(e)

		latency := time.Since(start)

		status := responseStatus(e, err)

		// ignore 404s for paths outside the API prefix (e.g. favicon.ico)
		if status == http.StatusNotFound && !strings.HasPrefix(req.URL.Path, "/"+config.APIVersion) {
			return err
		}

		attrs := []any{
			slog.String("method", req.Method),
			slog.String("path", req.URL.Path),
			slog.Int("status", status),
			slog.String("params", allParams),
			slog.Int64("latency_ms", latency.Milliseconds()),
			slog.String("ip", e.RealIP()),
		}

		if len(bodyBytes) > 0 {
			var bodyMap map[string]interface{}

			err := json.Unmarshal(bodyBytes, &bodyMap)
			if err == nil {
				body, err := json.Marshal(bodyMap)
				if err == nil {
					attrs = append(attrs, slog.String("body", string(body)))
				}
			}
		}

		if err != nil {
			attrs = append(attrs, slog.Any("error", err))
			slog.ErrorContext(e.Request().Context(), "request", attrs...)
		} else {
			slog.DebugContext(e.Request().Context(), "request", attrs...)
		}

		return err
	}
}

func getRequestParams(c echo.Context) string {
	// start with query parameters (e.g. ?someId=123)
	params := c.Request().URL.Query().Encode()

	// add path parameters (e.g. /:someId)
	for i, name := range c.ParamNames() {
		if params != "" {
			params += ", "
		}
		params += fmt.Sprintf("%s=%s", name, c.ParamValues()[i])
	}

	return params
}

// echo writes response after middleware runs, so c.Response().Status defaults to 200
// when an error is present the status must be read from the error instead
func responseStatus(c echo.Context, err error) int {
	if err != nil {
		var httpErr *echo.HTTPError
		if errors.As(err, &httpErr) {
			return httpErr.Code
		}
		return http.StatusInternalServerError
	}
	return c.Response().Status
}
