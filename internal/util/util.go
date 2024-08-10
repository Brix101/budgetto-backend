package util

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func NewRedisQueueClient(ctx context.Context, maxConns int) (*redis.Client, error) {
	return newRedisClient(ctx, "REDIS_URL", maxConns)
}

func newRedisClient(ctx context.Context, env string, maxConns int) (*redis.Client, error) {
	opt, err := redis.ParseURL(os.Getenv(env))
	if err != nil {
		return nil, err
	}
	opt.PoolSize = maxConns

	client := redis.NewClient(opt)

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return client, nil
}

func NewDatabasePool(ctx context.Context, maxConns int) (*pgxpool.Pool, error) {
	if maxConns == 0 {
		maxConns = 1
	}

	url := fmt.Sprintf(
		"%s?pool_max_conns=%d&pool_min_conns=%d",
		os.Getenv("DATABASE_URL"),
		maxConns,
		2,
	)
	config, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, err
	}

	// Setting the build statement cache to nil helps this work with pgbouncer
	config.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol
	config.MaxConnLifetime = 1 * time.Hour
	config.MaxConnIdleTime = 30 * time.Second
	return pgxpool.NewWithConfig(ctx, config)
}

func NewLogger(service string) *zap.Logger {
	env := os.Getenv("ENV")

	logger, _ := zap.NewProduction(zap.Fields(
		zap.String("env", env),
		zap.String("service", service),
	))

	if env == "" || env == "development" {
		logger, _ = zap.NewDevelopment()
	}

	return logger
}
