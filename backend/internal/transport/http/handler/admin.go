package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/anr-2609/anr-platform/backend/internal/platform/device"
	"github.com/anr-2609/anr-platform/backend/internal/platform/user"
	"github.com/anr-2609/anr-platform/backend/internal/transport/http/response"
)

type AdminLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AdminLoginResponse struct {
	AccessToken  string     `json:"access_token"`
	RefreshToken string     `json:"refresh_token"`
	TokenType    string     `json:"token_type"`
	User         *user.User `json:"user"`
}

type OverviewResponse struct {
	AppsCount    int64  `json:"apps_count"`
	DevicesCount int64  `json:"devices_count"`
	UsersCount   int64  `json:"users_count"`
	Status       string `json:"status"`
}

type ApplicationItem struct {
	AppID       string `json:"app_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

type AdminHandler struct {
	userService *user.Service
	userRepo    user.Repository
	deviceRepo  device.Repository
}

func NewAdminHandler(userService *user.Service, userRepo user.Repository, deviceRepo device.Repository) *AdminHandler {
	return &AdminHandler{
		userService: userService,
		userRepo:    userRepo,
		deviceRepo:  deviceRepo,
	}
}

func (h *AdminHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req AdminLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_REQUEST_BODY", "Failed to parse JSON body")
		return
	}

	u, accessToken, refreshToken, err := h.userService.Authenticate(r.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, user.ErrInvalidCredentials) {
			response.Error(w, http.StatusUnauthorized, "INVALID_CREDENTIALS", "Invalid email or password")
			return
		}
		if errors.Is(err, user.ErrAccountInactive) {
			response.Error(w, http.StatusForbidden, "ACCOUNT_INACTIVE", "Account is inactive or suspended")
			return
		}
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Authentication failed")
		return
	}

	response.JSON(w, http.StatusOK, AdminLoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		User:         u,
	})
}

func (h *AdminHandler) Overview(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	devicesCount, err := h.deviceRepo.Count(ctx)
	if err != nil {
		devicesCount = 0
	}

	usersCount, err := h.userRepo.Count(ctx)
	if err != nil {
		usersCount = 0
	}

	response.JSON(w, http.StatusOK, OverviewResponse{
		AppsCount:    1,
		DevicesCount: devicesCount,
		UsersCount:   usersCount,
		Status:       "healthy",
	})
}

func (h *AdminHandler) ListDevices(w http.ResponseWriter, r *http.Request) {
	devices, err := h.deviceRepo.List(r.Context(), 50, 0)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to list devices")
		return
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"devices": devices,
		"count":   len(devices),
	})
}

func (h *AdminHandler) ListApps(w http.ResponseWriter, r *http.Request) {
	apps := []ApplicationItem{
		{
			AppID:       "anr-001-wallpaper",
			Name:        "ANR Wallpaper",
			Description: "ANR Wallpaper Application",
			Status:      "active",
		},
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"apps":  apps,
		"count": len(apps),
	})
}
