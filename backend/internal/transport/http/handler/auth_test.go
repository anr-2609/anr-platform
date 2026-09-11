package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/anr-2609/anr-platform/backend/internal/platform/auth"
	"github.com/anr-2609/anr-platform/backend/internal/platform/device"
)

func setupTestAuthHandler() *AuthHandler {
	repo := device.NewMemoryRepository()
	devSvc := device.NewService(repo)
	tokSvc := auth.NewTokenService("test-secret-at-least-32-characters-long", 15*time.Minute, 24*time.Hour)
	return NewAuthHandler(devSvc, tokSvc)
}

func TestAuthHandler_DeviceSession(t *testing.T) {
	h := setupTestAuthHandler()

	body, _ := json.Marshal(DeviceSessionRequest{
		DeviceID:   "d-test-1",
		AppID:      "anr-001-wallpaper",
		Platform:   "android",
		OSVersion:  "14",
		AppVersion: "1.0.0",
	})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/device-session", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	h.DeviceSession(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var resp struct {
		Success bool                  `json:"success"`
		Data    DeviceSessionResponse `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if !resp.Success {
		t.Fatal("expected success true")
	}
	if resp.Data.AccessToken == "" || resp.Data.RefreshToken == "" {
		t.Fatal("expected access and refresh tokens")
	}
	if resp.Data.Device == nil || resp.Data.Device.DeviceID != "d-test-1" {
		t.Errorf("expected device ID d-test-1, got %v", resp.Data.Device)
	}
}

func TestAuthHandler_Refresh(t *testing.T) {
	h := setupTestAuthHandler()

	body, _ := json.Marshal(DeviceSessionRequest{
		DeviceID: "d-test-2",
		AppID:    "anr-001-wallpaper",
		Platform: "ios",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/device-session", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	h.DeviceSession(rec, req)

	var sessionResp struct {
		Data DeviceSessionResponse `json:"data"`
	}
	_ = json.Unmarshal(rec.Body.Bytes(), &sessionResp)

	refreshBody, _ := json.Marshal(RefreshRequest{
		RefreshToken: sessionResp.Data.RefreshToken,
	})
	refreshReq := httptest.NewRequest(http.MethodPost, "/api/v1/auth/refresh", bytes.NewReader(refreshBody))
	refreshRec := httptest.NewRecorder()

	h.Refresh(refreshRec, refreshReq)

	if refreshRec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", refreshRec.Code, refreshRec.Body.String())
	}

	var refreshResult struct {
		Success bool            `json:"success"`
		Data    RefreshResponse `json:"data"`
	}
	_ = json.Unmarshal(refreshRec.Body.Bytes(), &refreshResult)

	if !refreshResult.Success || refreshResult.Data.AccessToken == "" {
		t.Fatal("expected new access token")
	}
}
