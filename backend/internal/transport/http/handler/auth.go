package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/anr-2609/anr-platform/backend/internal/platform/auth"
	"github.com/anr-2609/anr-platform/backend/internal/platform/device"
	"github.com/anr-2609/anr-platform/backend/internal/transport/http/response"
)

type AuthHandler struct {
	deviceService *device.Service
	tokenService  *auth.TokenService
}

func NewAuthHandler(deviceService *device.Service, tokenService *auth.TokenService) *AuthHandler {
	return &AuthHandler{
		deviceService: deviceService,
		tokenService:  tokenService,
	}
}

type DeviceSessionRequest struct {
	DeviceID   string `json:"device_id"`
	AppID      string `json:"app_id"`
	Platform   string `json:"platform"`
	OSVersion  string `json:"os_version"`
	AppVersion string `json:"app_version"`
	PushToken  string `json:"push_token"`
}

type DeviceSessionResponse struct {
	AccessToken  string         `json:"access_token"`
	RefreshToken string         `json:"refresh_token"`
	TokenType    string         `json:"token_type"`
	Device       *device.Device `json:"device"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
}

func (h *AuthHandler) DeviceSession(w http.ResponseWriter, r *http.Request) {
	var req DeviceSessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid JSON payload")
		return
	}

	d, err := h.deviceService.RegisterOrUpdate(r.Context(), device.RegisterInput{
		DeviceID:   req.DeviceID,
		AppID:      req.AppID,
		Platform:   req.Platform,
		OSVersion:  req.OSVersion,
		AppVersion: req.AppVersion,
		PushToken:  req.PushToken,
	})
	if err != nil {
		if errors.Is(err, device.ErrInvalidDeviceID) || errors.Is(err, device.ErrInvalidAppID) || errors.Is(err, device.ErrInvalidPlatform) {
			response.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
			return
		}
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to register device")
		return
	}

	accessToken, refreshToken, err := h.tokenService.GenerateTokens(d.DeviceID, d.AppID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to generate tokens")
		return
	}

	response.JSON(w, http.StatusOK, DeviceSessionResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		Device:       d,
	})
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid JSON payload")
		return
	}

	claims, err := h.tokenService.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid or expired refresh token")
		return
	}

	accessToken, refreshToken, err := h.tokenService.GenerateTokens(claims.DeviceID, claims.AppID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to generate tokens")
		return
	}

	response.JSON(w, http.StatusOK, RefreshResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
	})
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.GetDeviceClaims(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "Missing credentials")
		return
	}

	d, err := h.deviceService.GetDevice(r.Context(), claims.DeviceID, claims.AppID)
	if err != nil {
		if errors.Is(err, device.ErrDeviceNotFound) {
			response.Error(w, http.StatusNotFound, "NOT_FOUND", "Device not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to retrieve device")
		return
	}

	response.JSON(w, http.StatusOK, d)
}
