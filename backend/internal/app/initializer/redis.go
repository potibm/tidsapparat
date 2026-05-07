package initializer

import (
	"context"
	"log/slog"
	"time"

	"github.com/potibm/billedapparat/internal/app/config"
	"github.com/redis/go-redis/v9"
)

const defaultRedisTimeout = 2 * time.Second

func InitializeRedis(dsn config.RedisURL) *redis.Client {
	logger := slog.With("component", "initializer.redis")

	if dsn == "" {
		logger.Info("No Redis DSN provided, skipping Redis initialization")

		return nil
	}

	if !dsn.IsValid() {
		logger.Error("Invalid Redis DSN", "dsn", dsn)

		return nil
	}

	options := dsn.RedisOptions()
	if options == nil {
		logger.Error("Failed to parse Redis DSN", "dsn", dsn)

		return nil
	}

	client := redis.NewClient(options)

	ctx, cancel := context.WithTimeout(context.Background(), defaultRedisTimeout)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		logger.Error("Redis enabled but connection failed", "dsn", dsn, "error", err)

		return nil
	}

	logger.Info("Successfully initialized Redis client")

	return client
}
