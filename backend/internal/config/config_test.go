package config

import (
	"os"
	"testing"
	"time"
)

func TestLoad_Default(t *testing.T) {
	_ = os.Unsetenv("SERVER_PORT")
	_ = os.Unsetenv("SERVER_ENV")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("expected no error loading config, got: %v", err)
	}

	if cfg.Server.Port != "8080" {
		t.Errorf("expected default port 8080, got %s", cfg.Server.Port)
	}

	if cfg.Server.Environment != "development" {
		t.Errorf("expected default env development, got %s", cfg.Server.Environment)
	}

	if cfg.Server.ReadTimeout != 10*time.Second {
		t.Errorf("expected default read timeout 10s, got %v", cfg.Server.ReadTimeout)
	}
}

func TestLoad_CustomEnv(t *testing.T) {
	t.Setenv("SERVER_PORT", "9090")
	t.Setenv("SERVER_ENV", "staging")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if cfg.Server.Port != "9090" {
		t.Errorf("expected port 9090, got %s", cfg.Server.Port)
	}

	if cfg.Server.Environment != "staging" {
		t.Errorf("expected env staging, got %s", cfg.Server.Environment)
	}
}

func TestDatabaseConfig_DSN(t *testing.T) {
	dbCfg := DatabaseConfig{
		Host:     "127.0.0.1",
		Port:     5432,
		User:     "postgres",
		Password: "secretpassword",
		Name:     "mydb",
		SSLMode:  "disable",
	}

	expected := "postgres://postgres:secretpassword@127.0.0.1:5432/mydb?sslmode=disable"
	if dbCfg.DSN() != expected {
		t.Errorf("expected DSN %s, got %s", expected, dbCfg.DSN())
	}
}
