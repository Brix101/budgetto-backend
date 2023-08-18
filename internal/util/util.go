package util

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewDatabasePool(ctx context.Context, database_url string, maxConns int) *pgxpool.Pool {
	if maxConns == 0 {
		maxConns = 1
	}

	url := fmt.Sprintf(
		"%s&pool_max_conns=%d&pool_min_conns=%d",
		database_url,
		maxConns,
		2,
	)
	config, err := pgxpool.ParseConfig(url)
	if err != nil {
		log.Fatal(err)
	}

	// Setting the build statement cache to nil helps this work with pgbouncer
	config.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol
	config.MaxConnLifetime = 1 * time.Hour
	config.MaxConnIdleTime = 30 * time.Second

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Fatal(err)
	}

	return pool
}
