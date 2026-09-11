package auth

import (
	"testing"
	"time"
)

func TestTokenService_GenerateAndValidate(t *testing.T) {
	svc := NewTokenService("test-secret-key-at-least-32-chars-long", 15*time.Minute, 24*time.Hour)

	accessToken, refreshToken, err := svc.GenerateTokens("device-123", "anr-001-wallpaper")
	if err != nil {
		t.Fatalf("unexpected error generating tokens: %v", err)
	}

	if accessToken == "" || refreshToken == "" {
		t.Fatal("expected non-empty tokens")
	}

	accessClaims, err := svc.ValidateAccessToken(accessToken)
	if err != nil {
		t.Fatalf("unexpected error validating access token: %v", err)
	}

	if accessClaims.DeviceID != "device-123" {
		t.Errorf("expected device-123, got %s", accessClaims.DeviceID)
	}
	if accessClaims.AppID != "anr-001-wallpaper" {
		t.Errorf("expected anr-001-wallpaper, got %s", accessClaims.AppID)
	}
	if accessClaims.Type != "GUEST_ACCESS" {
		t.Errorf("expected GUEST_ACCESS, got %s", accessClaims.Type)
	}

	refreshClaims, err := svc.ValidateRefreshToken(refreshToken)
	if err != nil {
		t.Fatalf("unexpected error validating refresh token: %v", err)
	}

	if refreshClaims.Type != "GUEST_REFRESH" {
		t.Errorf("expected GUEST_REFRESH, got %s", refreshClaims.Type)
	}

	_, err = svc.ValidateAccessToken(refreshToken)
	if err == nil {
		t.Fatal("expected error when validating refresh token as access token")
	}
}

func TestTokenService_ExpiredToken(t *testing.T) {
	svc := NewTokenService("test-secret-key-at-least-32-chars-long", -1*time.Minute, -1*time.Minute)

	accessToken, _, err := svc.GenerateTokens("device-123", "anr-001-wallpaper")
	if err != nil {
		t.Fatalf("unexpected error generating token: %v", err)
	}

	_, err = svc.ValidateAccessToken(accessToken)
	if err == nil {
		t.Fatal("expected error for expired token, got nil")
	}
}
