package database

import (
	"context"
	"fmt"
	"github.com/exaring/otelpgx"
	"github.com/jackc/pgx/v5/pgxpool"
	"guilhermefaleiros/go-service-template/internal/shared"
	"time"
)

func NewPGConnection(ctx context.Context, config *shared.Config) (*pgxpool.Pool, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		config.Postgres.User, config.Postgres.Password, config.Postgres.Host, config.Postgres.Port, config.Postgres.Name)
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	poolConfig, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("unable to parse connection string: %w", err)
	}

	poolConfig.MaxConns = config.Postgres.MaxConnection
	poolConfig.MinConns = config.Postgres.MinConnection
	poolConfig.MaxConnIdleTime = time.Duration(config.Postgres.MaxIdleTime) * time.Second
	poolConfig.MaxConnLifetime = time.Duration(config.Postgres.MaxLifeTime) * time.Second

	poolConfig.ConnConfig.Tracer = otelpgx.NewTracer()

	pool, err := pgxpool.NewWithConfig(ctxWithTimeout, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}
	return pool, nil
}
