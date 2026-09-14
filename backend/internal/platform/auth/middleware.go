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
	UserClaimsKey   contextKey = "user_claims"
)

func RequireDeviceAuth(tokenService *TokenService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenStr := extractBearerToken(r)
			if tokenStr == "" {
				response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "Valid bearer token is required")
				return
			}

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

func RequireUserAuth(tokenService *TokenService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenStr := extractBearerToken(r)
			if tokenStr == "" {
				response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "Valid bearer token is required")
				return
			}

			claims, err := tokenService.ValidateUserAccessToken(tokenStr)
			if err != nil {
				response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid or expired user access token")
				return
			}

			ctx := context.WithValue(r.Context(), UserClaimsKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequireRole(requiredRole string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := GetUserClaims(r.Context())
			if !ok || claims == nil {
				response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication required")
				return
			}

			if claims.Role != requiredRole {
				response.Error(w, http.StatusForbidden, "FORBIDDEN", "Insufficient permissions")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func GetUserClaims(ctx context.Context) (*UserClaims, bool) {
	if ctx == nil {
		return nil, false
	}
	claims, ok := ctx.Value(UserClaimsKey).(*UserClaims)
	return claims, ok
}

func extractBearerToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return ""
	}

	return strings.TrimSpace(parts[1])
}
