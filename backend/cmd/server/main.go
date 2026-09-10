package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/anr-2609/anr-platform/backend/internal/cache"
	"github.com/anr-2609/anr-platform/backend/internal/config"
	"github.com/anr-2609/anr-platform/backend/internal/database"
	transporthttp "github.com/anr-2609/anr-platform/backend/internal/transport/http"
	"github.com/anr-2609/anr-platform/backend/internal/transport/http/handler"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "application error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	// 1. Tải cấu hình môi trường
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("loading configuration: %w", err)
	}

	// 2. Khởi tạo Structured Logger (slog)
	var logHandler slog.Handler
	if cfg.Server.Environment == "production" {
		logHandler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	} else {
		logHandler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	}
	logger := slog.New(logHandler)
	slog.SetDefault(logger)

	logger.Info("starting ANR Platform Backend",
		slog.String("env", cfg.Server.Environment),
		slog.String("port", cfg.Server.Port),
	)

	// 3. Khởi tạo Database (nếu enabled)
	var db *database.Postgres
	if cfg.Database.ConnectionEnabled {
		logger.Info("connecting to PostgreSQL database...", slog.String("host", cfg.Database.Host))
		var dbErr error
		db, dbErr = database.NewPostgres(context.Background(), cfg.Database)
		if dbErr != nil {
			logger.Error("failed to connect to PostgreSQL", slog.String("error", dbErr.Error()))
			return fmt.Errorf("initializing postgres: %w", dbErr)
		}
		defer db.Close()
		logger.Info("connected to PostgreSQL successfully")
	} else {
		logger.Info("PostgreSQL connection disabled by configuration")
	}

	// 4. Khởi tạo Redis Cache (nếu enabled)
	var redisCache *cache.Redis
	if cfg.Redis.Enabled {
		logger.Info("connecting to Redis...", slog.String("host", cfg.Redis.Host))
		var redisErr error
		redisCache, redisErr = cache.NewRedis(context.Background(), cfg.Redis)
		if redisErr != nil {
			logger.Error("failed to connect to Redis", slog.String("error", redisErr.Error()))
			return fmt.Errorf("initializing redis: %w", redisErr)
		}
		defer func() { _ = redisCache.Close() }()
		logger.Info("connected to Redis successfully")
	} else {
		logger.Info("Redis connection disabled by configuration")
	}

	// 5. Khởi tạo Handlers (Manual Constructor Injection)
	var dbPinger handler.Pinger
	if db != nil {
		dbPinger = db
	}
	var cachePinger handler.Pinger
	if redisCache != nil {
		cachePinger = redisCache
	}
	healthHandler := handler.NewHealthHandler(dbPinger, cachePinger)

	// 6. Khởi tạo Router & Middleware
	router := transporthttp.NewRouter(transporthttp.RouterConfig{
		Logger:        logger,
		HealthHandler: healthHandler,
	})

	// 7. Khởi tạo HTTP Server
	server := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// 8. Lắng nghe tín hiệu Graceful Shutdown từ OS
	serverErrors := make(chan error, 1)
	go func() {
		logger.Info("server is listening", slog.String("address", server.Addr))
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErrors <- fmt.Errorf("server error: %w", err)
		}
	}()

	shutdownSignal := make(chan os.Signal, 1)
	signal.Notify(shutdownSignal, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)
	case sig := <-shutdownSignal:
		logger.Info("shutdown signal received", slog.String("signal", sig.String()))

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			logger.Error("server forced to shutdown", slog.String("error", err.Error()))
			_ = server.Close()
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}

		logger.Info("server stopped gracefully")
	}

	return nil
}
