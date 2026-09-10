package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/anr-2609/anr-platform/backend/internal/config"
)

// Redis đóng gói redis.Client.
type Redis struct {
	Client *redis.Client
}

// NewRedis khởi tạo kết nối tới Redis instance.
func NewRedis(ctx context.Context, cfg config.RedisConfig) (*Redis, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	pingCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if err := client.Ping(pingCtx).Err(); err != nil {
		_ = client.Close()
		return nil, fmt.Errorf("pinging redis: %w", err)
	}

	return &Redis{Client: client}, nil
}

// Ping kiểm tra kết nối với Redis.
func (r *Redis) Ping(ctx context.Context) error {
	if r == nil || r.Client == nil {
		return fmt.Errorf("redis client is uninitialized")
	}
	return r.Client.Ping(ctx).Err()
}

// Close ngắt kết nối với Redis.
func (r *Redis) Close() error {
	if r != nil && r.Client != nil {
		return r.Client.Close()
	}
	return nil
}
