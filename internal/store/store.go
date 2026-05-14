package store

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"

	productstore "github.com/keto-granola/server/internal/product/store"
	"github.com/keto-granola/server/internal/store/db/generated"
)

type Store struct {
	pool    *pgxpool.Pool
	queries *generated.Queries
}

func New(ctx context.Context) (*Store, error) {
	return &Store{
		pool:    nil,
		queries: nil,
	}, nil
}

func (s *Store) ProductStore() *productstore.Store {
	return productstore.New(s.queries)
}

func (s *Store) Close() {
	slog.Info("closing pool", slog.Any("pool", s.pool))
}
