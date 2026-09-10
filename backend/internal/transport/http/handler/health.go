package handler

import (
	"context"
	"net/http"

	"github.com/anr-2609/anr-platform/backend/internal/transport/http/response"
)

// Pinger là interface kiểm tra kết nối tới các external dependencies.
type Pinger interface {
	Ping(ctx context.Context) error
}

type HealthHandler struct {
	db    Pinger
	cache Pinger
}

func NewHealthHandler(db Pinger, cache Pinger) *HealthHandler {
	return &HealthHandler{
		db:    db,
		cache: cache,
	}
}

// Livez kiểm tra tiến trình có đang sống hay không (K8s / Docker healthcheck).
func (h *HealthHandler) Livez(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, map[string]string{
		"status": "alive",
	})
}

// Readyz kiểm tra hệ thống và các dịch vụ phụ thuộc (DB, Redis) đã sẵn sàng phục vụ traffic chưa.
func (h *HealthHandler) Readyz(w http.ResponseWriter, r *http.Request) {
	checks := make(map[string]string)
	isReady := true

	if h.db != nil {
		if err := h.db.Ping(r.Context()); err != nil {
			checks["database"] = "down: " + err.Error()
			isReady = false
		} else {
			checks["database"] = "up"
		}
	} else {
		checks["database"] = "disabled"
	}

	if h.cache != nil {
		if err := h.cache.Ping(r.Context()); err != nil {
			checks["cache"] = "down: " + err.Error()
			isReady = false
		} else {
			checks["cache"] = "up"
		}
	} else {
		checks["cache"] = "disabled"
	}

	if !isReady {
		response.Error(w, http.StatusServiceUnavailable, "SERVICE_UNAVAILABLE", "One or more dependent services are unavailable")
		return
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"status": "ready",
		"checks": checks,
	})
}

// Healthz trả về tổng quan trạng thái service.
func (h *HealthHandler) Healthz(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, map[string]interface{}{
		"service": "anr-platform-backend",
		"status":  "healthy",
		"version": "1.0.0",
	})
}
