package product

import (
	"context"
)

type Repository interface {
	Add(ctx context.Context, product *AddProductInput) (*Product, error)
}
