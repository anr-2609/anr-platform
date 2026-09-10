package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/anr-2609/anr-platform/backend/internal/transport/http/response"
)

// Recoverer bắt các panic chưa được xử lý, log stack trace và trả HTTP 500.
func Recoverer(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rvr := recover(); rvr != nil {
					if rvr == http.ErrAbortHandler {
						panic(rvr)
					}

					reqID := GetRequestID(r.Context())
					stack := string(debug.Stack())

					logger.ErrorContext(r.Context(), "panic recovered",
						slog.String("request_id", reqID),
						slog.String("panic", fmt.Sprintf("%v", rvr)),
						slog.String("stack", stack),
					)

					response.Error(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "An unexpected internal server error occurred")
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
