package server

import (
	"github.com/labstack/echo/v4"
)

func registerRoutes(public, private *echo.Group, handlers *Handlers) {
	public.POST("/health", func(c echo.Context) error { return nil })

	// admin routes
	private.POST("/admin/product", handlers.ProductAdmin.CreateProduct)
}
