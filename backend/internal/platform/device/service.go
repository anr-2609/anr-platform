package device

import (
	"context"
	"errors"
	"strings"
)

var (
	ErrInvalidDeviceID = errors.New("device_id is required")
	ErrInvalidAppID    = errors.New("app_id is required")
	ErrInvalidPlatform = errors.New("platform must be ios, android, or web")
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

type RegisterInput struct {
	DeviceID   string
	AppID      string
	Platform   string
	OSVersion  string
	AppVersion string
	PushToken  string
}

func (s *Service) RegisterOrUpdate(ctx context.Context, input RegisterInput) (*Device, error) {
	input.DeviceID = strings.TrimSpace(input.DeviceID)
	input.AppID = strings.TrimSpace(input.AppID)
	input.Platform = strings.ToLower(strings.TrimSpace(input.Platform))

	if input.DeviceID == "" {
		return nil, ErrInvalidDeviceID
	}
	if input.AppID == "" {
		return nil, ErrInvalidAppID
	}
	if input.Platform != "android" && input.Platform != "ios" && input.Platform != "web" {
		return nil, ErrInvalidPlatform
	}

	d := &Device{
		DeviceID:   input.DeviceID,
		AppID:      input.AppID,
		Platform:   input.Platform,
		OSVersion:  input.OSVersion,
		AppVersion: input.AppVersion,
		PushToken:  input.PushToken,
	}

	return s.repo.Upsert(ctx, d)
}

func (s *Service) GetDevice(ctx context.Context, deviceID, appID string) (*Device, error) {
	return s.repo.GetByDeviceAndApp(ctx, deviceID, appID)
}
