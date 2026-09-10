package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
)

type contextKey string

const (
	RequestIDHeader             = "X-Request-ID"
	RequestIDKey     contextKey = "request_id"
)

// RequestID middleware đảm bảo mỗi HTTP request đều có Request ID duy nhất.
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Header.Get(RequestIDHeader)
		if reqID == "" {
			reqID = generateID()
		}

		ctx := context.WithValue(r.Context(), RequestIDKey, reqID)
		w.Header().Set(RequestIDHeader, reqID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetRequestID lấy Request ID từ context nếu có.
func GetRequestID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if reqID, ok := ctx.Value(RequestIDKey).(string); ok {
		return reqID
	}
	return ""
}

func generateID() string {
	b := make([]byte, 12)
	_, err := rand.Read(b)
	if err != nil {
		return "unknown-req-id"
	}
	return hex.EncodeToString(b)
}
