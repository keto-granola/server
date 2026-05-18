package server

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/keto-granola/server/internal/store"
)

const pingTimeout = 5 * time.Second

func registerRoutes(public, private *echo.Group, handlers *Handlers, dataStore *store.Store) {
	registerHealthEndpoint(public, dataStore)

	// admin routes
	private.POST("/admin/product", Handle(handlers.ProductAdmin.CreateProduct, http.StatusCreated))
}

func registerHealthEndpoint(public *echo.Group, dataStore *store.Store) {
	public.POST("/health", func(e echo.Context) error {
		dbStatus := "ok"
		httpStatus := http.StatusOK

		pingCtx, cancel := context.WithTimeout(e.Request().Context(), pingTimeout)
		defer cancel()

		err := dataStore.PingDB(pingCtx)
		if err != nil {
			httpStatus = http.StatusServiceUnavailable
			dbStatus = "unreachable"
		}

		return e.JSON(httpStatus, map[string]string{
			"status": "ok",
			"db":     dbStatus,
		})
	})
}
