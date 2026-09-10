package http

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/anr-2609/anr-platform/backend/internal/transport/http/handler"
	"github.com/anr-2609/anr-platform/backend/internal/transport/http/middleware"
	"github.com/anr-2609/anr-platform/backend/internal/transport/http/response"
)

// RouterConfig chứa các handlers và dependencies cần thiết để dựng HTTP router.
type RouterConfig struct {
	Logger        *slog.Logger
	HealthHandler *handler.HealthHandler
}

// NewRouter khởi tạo chi router với đầy đủ middleware và route groupings.
func NewRouter(cfg RouterConfig) http.Handler {
	r := chi.NewRouter()

	// Global Middlewares
	r.Use(chimw.RealIP)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger(cfg.Logger))
	r.Use(middleware.Recoverer(cfg.Logger))

	// CORS Setup
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*.anr-studio.com", "http://localhost:*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Request-ID", "X-App-ID", "X-Device-ID"},
		ExposedHeaders:   []string{"Link", "X-Request-ID"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Operational & Health Endpoints
	r.Get("/livez", cfg.HealthHandler.Livez)
	r.Get("/readyz", cfg.HealthHandler.Readyz)
	r.Get("/healthz", cfg.HealthHandler.Healthz)

	// API V1 Group
	r.Route("/api/v1", func(v1 chi.Router) {
		// Common Platform APIs
		v1.Route("/platform", func(platform chi.Router) {
			platform.Get("/info", func(w http.ResponseWriter, r *http.Request) {
				response.JSON(w, http.StatusOK, map[string]string{
					"name":    "ANR Platform Core",
					"version": "1.0.0",
				})
			})
		})

		// Admin APIs (Protected)
		v1.Route("/admin", func(admin chi.Router) {
			admin.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
				response.JSON(w, http.StatusOK, map[string]string{
					"message": "admin endpoint placeholder",
				})
			})
		})

		// App Modules APIs (anr-NNN-slug)
		v1.Route("/modules", func(modules chi.Router) {
			modules.Get("/catalog", func(w http.ResponseWriter, r *http.Request) {
				response.JSON(w, http.StatusOK, map[string]interface{}{
					"modules": []string{"anr-001-wallpaper"},
				})
			})
		})
	})

	// 404 handler chuẩn hoá JSON
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		response.Error(w, http.StatusNotFound, "NOT_FOUND", "The requested resource was not found")
	})

	// 405 handler
	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		response.Error(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "HTTP method is not allowed")
	})

	return r
}
