package initializer

import (
	"context"
	"log/slog"
	"time"

	"github.com/potibm/billedapparat/internal/app/config"
	"github.com/redis/go-redis/v9"
)

func InitializeRedis(dsn config.RedisURL) *redis.Client {
	slog := slog.With("component", "initializer.redis")

	if dsn == "" {
		slog.Info("No Redis DSN provided, skipping Redis initialization")

		return nil
	}

	if !dsn.IsValid() {
		slog.Error("Invalid Redis DSN", "dsn", dsn)

		return nil
	}

	options := dsn.RedisOptions()
	if options == nil {
		slog.Error("Failed to parse Redis DSN", "dsn", dsn)

		return nil
	}

	client := redis.NewClient(options)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		slog.Error("Redis enabled but connection failed", "dsn", dsn, "error", err)

		return nil
	}

	slog.Info("Successfully initialized Redis client")

	return client
}
