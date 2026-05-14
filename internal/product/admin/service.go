package admin

import (
	"context"

	"github.com/keto-granola/server/internal/product"
)

type Service struct {
	store product.Repository
}

func NewService(store product.Repository) *Service {
	return &Service{store: store}
}

func (s *Service) CreateProduct(ctx context.Context, req CreateProductRequest) (*product.Product, error) {
	p := &product.Product{
		Name: req.Name,
	}

	return s.store.InsertProduct(ctx, p)
}
