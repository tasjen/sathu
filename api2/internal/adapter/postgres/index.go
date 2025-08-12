package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tasjen/sathu/api-hexa/internal/adapter/config"
	sqlc_gen "github.com/tasjen/sathu/api-hexa/internal/adapter/postgres/sqlc/gen"
)

type DB struct {
	Pool    *pgxpool.Pool
	Queries *sqlc_gen.Queries
}

func New(ctx context.Context, config *config.DB) (*DB, error) {
	uri := fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=disable",
		config.Connection,
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Name,
	)
	pool, err := pgxpool.New(ctx, uri)
	if err != nil {
		return nil, err
	}
	err = pool.Ping(ctx)
	if err != nil {
		return nil, err
	}
	queries := sqlc_gen.New(pool)
	return &DB{Pool: pool, Queries: queries}, nil
}

func ToPgError(err error) (*pgconn.PgError, bool) {
	pgErr, ok := err.(*pgconn.PgError)
	if !ok {
		return nil, false
	}
	return pgErr, true
}
