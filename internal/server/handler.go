package server

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler[Req any, Res any] func(context.Context, Req) (Res, error)

func Handle[Req any, Res any](handler Handler[Req, Res], statusCode int) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var in Req

		if err := ctx.Bind(&in); err != nil {
			return echo.NewHTTPError(
				http.StatusBadRequest,
				"invalid request",
			)
		}

		out, err := handler(ctx.Request().Context(), in)
		if err != nil {
			return echo.NewHTTPError(
				http.StatusInternalServerError,
				"internal server error",
			).SetInternal(err)
		}

		return ctx.JSON(statusCode, out)
	}
}
