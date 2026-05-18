package admin

import (
	"context"

	"github.com/keto-granola/server/internal/product"
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

func (h *Handler) CreateProduct(ctx context.Context, req CreateProductRequest) (*product.Product, error) {
	prod, err := h.service.CreateProduct(ctx, req)
	if err != nil {
		return nil, err
	}

	return prod, nil
}
