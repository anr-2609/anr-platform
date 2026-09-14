package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/anr-2609/anr-platform/backend/internal/platform/auth"
	"github.com/anr-2609/anr-platform/backend/internal/platform/device"
	"github.com/anr-2609/anr-platform/backend/internal/platform/user"
)

func TestAdminHandler_Login(t *testing.T) {
	ctx := context.Background()
	userRepo := user.NewMemoryRepository()
	deviceRepo := device.NewMemoryRepository()
	tokenService := auth.NewTokenService("test_secret_admin_handler_32_chars_long", 15*time.Minute, 7*24*time.Hour)
	userService := user.NewService(userRepo, tokenService)

	err := userService.EnsureAdmin(ctx, "admin@anr-studio.com", "AdminSecret123!")
	require.NoError(t, err)

	adminHandler := NewAdminHandler(userService, userRepo, deviceRepo)

	t.Run("Success", func(t *testing.T) {
		body, _ := json.Marshal(AdminLoginRequest{
			Email:    "admin@anr-studio.com",
			Password: "AdminSecret123!",
		})

		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(body))
		rec := httptest.NewRecorder()

		adminHandler.Login(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var resp struct {
			Success bool               `json:"success"`
			Data    AdminLoginResponse `json:"data"`
		}
		err := json.NewDecoder(rec.Body).Decode(&resp)
		require.NoError(t, err)
		assert.True(t, resp.Success)
		assert.NotEmpty(t, resp.Data.AccessToken)
		assert.Equal(t, "admin@anr-studio.com", resp.Data.User.Email)
	})

	t.Run("InvalidCredentials", func(t *testing.T) {
		body, _ := json.Marshal(AdminLoginRequest{
			Email:    "admin@anr-studio.com",
			Password: "WrongPassword!",
		})

		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(body))
		rec := httptest.NewRecorder()

		adminHandler.Login(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})
}

func TestAdminHandler_Overview_And_Devices(t *testing.T) {
	ctx := context.Background()
	userRepo := user.NewMemoryRepository()
	deviceRepo := device.NewMemoryRepository()
	tokenService := auth.NewTokenService("test_secret_admin_handler_32_chars_long", 15*time.Minute, 7*24*time.Hour)
	userService := user.NewService(userRepo, tokenService)

	_, err := deviceRepo.Upsert(ctx, &device.Device{
		DeviceID: "device-1",
		AppID:    "anr-001-wallpaper",
		Platform: "android",
	})
	require.NoError(t, err)

	adminHandler := NewAdminHandler(userService, userRepo, deviceRepo)

	t.Run("Overview", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/overview", nil)
		rec := httptest.NewRecorder()

		adminHandler.Overview(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var resp struct {
			Success bool             `json:"success"`
			Data    OverviewResponse `json:"data"`
		}
		err := json.NewDecoder(rec.Body).Decode(&resp)
		require.NoError(t, err)
		assert.True(t, resp.Success)
		assert.Equal(t, int64(1), resp.Data.DevicesCount)
		assert.Equal(t, int64(1), resp.Data.AppsCount)
	})

	t.Run("ListDevices", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/devices", nil)
		rec := httptest.NewRecorder()

		adminHandler.ListDevices(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var resp struct {
			Success bool `json:"success"`
			Data    struct {
				Devices []*device.Device `json:"devices"`
				Count   int              `json:"count"`
			} `json:"data"`
		}
		err := json.NewDecoder(rec.Body).Decode(&resp)
		require.NoError(t, err)
		assert.True(t, resp.Success)
		assert.Equal(t, 1, resp.Data.Count)
		assert.Equal(t, "device-1", resp.Data.Devices[0].DeviceID)
	})
}
