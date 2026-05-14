package product

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	InsertProduct(ctx context.Context, product *Product) (*Product, error)
}

type Product struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
