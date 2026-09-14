package http

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/anr-2609/anr-platform/backend/internal/platform/auth"
	"github.com/anr-2609/anr-platform/backend/internal/transport/http/handler"
	"github.com/anr-2609/anr-platform/backend/internal/transport/http/middleware"
	"github.com/anr-2609/anr-platform/backend/internal/transport/http/response"
)

type RouterConfig struct {
	Logger        *slog.Logger
	HealthHandler *handler.HealthHandler
	AuthHandler   *handler.AuthHandler
	AdminHandler  *handler.AdminHandler
	TokenService  *auth.TokenService
}

func NewRouter(cfg RouterConfig) http.Handler {
	r := chi.NewRouter()

	r.Use(chimw.RealIP)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger(cfg.Logger))
	r.Use(middleware.Recoverer(cfg.Logger))

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*.anr-studio.com", "http://localhost:*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Request-ID", "X-App-ID", "X-Device-ID"},
		ExposedHeaders:   []string{"Link", "X-Request-ID"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/livez", cfg.HealthHandler.Livez)
	r.Get("/readyz", cfg.HealthHandler.Readyz)
	r.Get("/healthz", cfg.HealthHandler.Healthz)

	r.Route("/api/v1", func(v1 chi.Router) {
		v1.Route("/auth", func(authRouter chi.Router) {
			authRouter.Post("/device-session", cfg.AuthHandler.DeviceSession)
			authRouter.Post("/refresh", cfg.AuthHandler.Refresh)
			if cfg.AdminHandler != nil {
				authRouter.Post("/login", cfg.AdminHandler.Login)
			}
		})

		v1.Group(func(protected chi.Router) {
			protected.Use(auth.RequireDeviceAuth(cfg.TokenService))

			protected.Get("/devices/me", cfg.AuthHandler.Me)

			protected.Route("/platform", func(platform chi.Router) {
				platform.Get("/info", func(w http.ResponseWriter, r *http.Request) {
					claims, _ := auth.GetDeviceClaims(r.Context())
					response.JSON(w, http.StatusOK, map[string]interface{}{
						"name":      "ANR Platform Core",
						"version":   "1.0.0",
						"device_id": claims.DeviceID,
						"app_id":    claims.AppID,
					})
				})
			})

			protected.Route("/modules", func(modules chi.Router) {
				modules.Get("/catalog", func(w http.ResponseWriter, r *http.Request) {
					response.JSON(w, http.StatusOK, map[string]interface{}{
						"modules": []string{"anr-001-wallpaper"},
					})
				})
			})
		})

		if cfg.AdminHandler != nil {
			v1.Route("/admin", func(admin chi.Router) {
				admin.Use(auth.RequireUserAuth(cfg.TokenService))
				admin.Use(auth.RequireRole("admin"))

				admin.Get("/overview", cfg.AdminHandler.Overview)
				admin.Get("/apps", cfg.AdminHandler.ListApps)
				admin.Get("/devices", cfg.AdminHandler.ListDevices)
			})
		}
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		response.Error(w, http.StatusNotFound, "NOT_FOUND", "The requested resource was not found")
	})

	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		response.Error(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "HTTP method is not allowed")
	})

	return r
}
