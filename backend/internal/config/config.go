package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config chứa toàn bộ cấu hình ứng dụng backend.
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	Auth     AuthConfig
}

type ServerConfig struct {
	Port         string
	Environment  string // "development" | "production"
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

type DatabaseConfig struct {
	Host               string
	Port               int
	User               string
	Password           string
	Name               string
	SSLMode            string
	MaxConns           int32
	MinConns           int32
	MaxConnLifetime    time.Duration
	MaxConnIdleTime    time.Duration
	ConnectionEnabled  bool
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
	Enabled  bool
}

type AuthConfig struct {
	JWTSecret        string
	AccessTokenTTL   time.Duration
	RefreshTokenTTL  time.Duration
}

// Load đọc cấu hình từ biến môi trường và gán giá trị mặc định an toàn.
func Load() (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			Port:         getEnv("SERVER_PORT", "8080"),
			Environment:  getEnv("SERVER_ENV", "development"),
			ReadTimeout:  getEnvDuration("SERVER_TIMEOUT_READ", 10*time.Second),
			WriteTimeout: getEnvDuration("SERVER_TIMEOUT_WRITE", 10*time.Second),
			IdleTimeout:  getEnvDuration("SERVER_TIMEOUT_IDLE", 120*time.Second),
		},
		Database: DatabaseConfig{
			Host:              getEnv("DB_HOST", "localhost"),
			Port:              getEnvInt("DB_PORT", 5432),
			User:              getEnv("DB_USER", "anr_user"),
			Password:          getEnv("DB_PASSWORD", "anr_password_change_me"),
			Name:              getEnv("DB_NAME", "anr_platform"),
			SSLMode:           getEnv("DB_SSLMODE", "disable"),
			MaxConns:          int32(getEnvInt("DB_MAX_CONNS", 25)),
			MinConns:          int32(getEnvInt("DB_MIN_CONNS", 5)),
			MaxConnLifetime:   getEnvDuration("DB_CONN_LIFETIME", 1*time.Hour),
			MaxConnIdleTime:   getEnvDuration("DB_CONN_IDLE_TIME", 30*time.Minute),
			ConnectionEnabled: getEnvBool("DB_ENABLED", false), // Mặc định tắt để test standalone
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnvInt("REDIS_PORT", 6379),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvInt("REDIS_DB", 0),
			Enabled:  getEnvBool("REDIS_ENABLED", false), // Mặc định tắt để test standalone
		},
		Auth: AuthConfig{
			JWTSecret:       getEnv("JWT_SECRET", "change-me-in-production-at-least-32-chars-long"),
			AccessTokenTTL:  getEnvDuration("JWT_ACCESS_TTL", 15*time.Minute),
			RefreshTokenTTL: getEnvDuration("JWT_REFRESH_TTL", 7*24*time.Hour),
		},
	}

	if cfg.Server.Environment == "production" && cfg.Auth.JWTSecret == "change-me-in-production-at-least-32-chars-long" {
		return nil, fmt.Errorf("in production, JWT_SECRET must be explicitly configured")
	}

	return cfg, nil
}

// DSN trả về PostgreSQL connection string định dạng URL.
func (d *DatabaseConfig) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		d.User, d.Password, d.Host, d.Port, d.Name, d.SSLMode)
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	valStr := getEnv(key, "")
	if valStr == "" {
		return fallback
	}
	val, err := strconv.Atoi(valStr)
	if err != nil {
		return fallback
	}
	return val
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	valStr := getEnv(key, "")
	if valStr == "" {
		return fallback
	}
	val, err := time.ParseDuration(valStr)
	if err != nil {
		return fallback
	}
	return val
}

func getEnvBool(key string, fallback bool) bool {
	valStr := getEnv(key, "")
	if valStr == "" {
		return fallback
	}
	val, err := strconv.ParseBool(valStr)
	if err != nil {
		return fallback
	}
	return val
}
