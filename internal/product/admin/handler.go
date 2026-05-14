package admin

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{service: s}
}

type CreateProductRequest struct {
	Name string `json:"name" validate:"required"`
}

func (h *Handler) CreateProduct(e echo.Context) error {
	ctx := e.Request().Context()

	var req CreateProductRequest

	product, err := h.service.CreateProduct(ctx, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "error creating new product").SetInternal(err)
	}

	return e.JSON(http.StatusCreated, product)
}
