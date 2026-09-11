package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/anr-2609/anr-platform/backend/internal/transport/http/response"
)

type contextKey string

const (
	DeviceClaimsKey contextKey = "device_claims"
)

func RequireDeviceAuth(tokenService *TokenService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authorization header is required")
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid authorization header format")
				return
			}

			tokenStr := strings.TrimSpace(parts[1])
			claims, err := tokenService.ValidateAccessToken(tokenStr)
			if err != nil {
				response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid or expired access token")
				return
			}

			ctx := context.WithValue(r.Context(), DeviceClaimsKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetDeviceClaims(ctx context.Context) (*DeviceClaims, bool) {
	if ctx == nil {
		return nil, false
	}
	claims, ok := ctx.Value(DeviceClaimsKey).(*DeviceClaims)
	return claims, ok
}
