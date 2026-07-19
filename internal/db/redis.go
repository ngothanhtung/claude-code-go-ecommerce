package db

import (
	"context"

	"github.com/redis/go-redis/v9"

	"github.com/ngothanhtung/go-tutorials/internal/config"
)

// NewRedis opens a redis client.
func NewRedis(cfg config.RedisConfig) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
}

// PingRedis checks the connection is alive.
func PingRedis(ctx context.Context, rdb *redis.Client) error {
	return rdb.Ping(ctx).Err()
}
