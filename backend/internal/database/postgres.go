package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/anr-2609/anr-platform/backend/internal/config"
)

// Postgres đóng gói pgxpool.Pool và cung cấp các hàm quản trị connection pool.
type Postgres struct {
	Pool *pgxpool.Pool
}

// NewPostgres khởi tạo connection pool tới PostgreSQL với cấu hình tối ưu.
func NewPostgres(ctx context.Context, cfg config.DatabaseConfig) (*Postgres, error) {
	poolConfig, err := pgxpool.ParseConfig(cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("parsing postgres connection string: %w", err)
	}

	poolConfig.MaxConns = cfg.MaxConns
	poolConfig.MinConns = cfg.MinConns
	poolConfig.MaxConnLifetime = cfg.MaxConnLifetime
	poolConfig.MaxConnIdleTime = cfg.MaxConnIdleTime

	// Connection timeout khi khởi tạo
	connectCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(connectCtx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("creating postgres connection pool: %w", err)
	}

	if err := pool.Ping(connectCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("pinging postgres database: %w", err)
	}

	return &Postgres{Pool: pool}, nil
}

// Ping kiểm tra kết nối với PostgreSQL.
func (p *Postgres) Ping(ctx context.Context) error {
	if p == nil || p.Pool == nil {
		return fmt.Errorf("postgres pool is uninitialized")
	}
	return p.Pool.Ping(ctx)
}

// Close đóng toàn bộ connection pool.
func (p *Postgres) Close() {
	if p != nil && p.Pool != nil {
		p.Pool.Close()
	}
}
