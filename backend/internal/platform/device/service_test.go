package device

import (
	"context"
	"errors"
	"testing"
)

func TestService_RegisterOrUpdate(t *testing.T) {
	repo := NewMemoryRepository()
	svc := NewService(repo)
	ctx := context.Background()

	d, err := svc.RegisterOrUpdate(ctx, RegisterInput{
		DeviceID:   "dev-abc-1",
		AppID:      "anr-001-wallpaper",
		Platform:   "android",
		OSVersion:  "14",
		AppVersion: "1.0.0",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if d.ID != 1 {
		t.Errorf("expected ID 1, got %d", d.ID)
	}
	if d.DeviceID != "dev-abc-1" {
		t.Errorf("expected DeviceID dev-abc-1, got %s", d.DeviceID)
	}

	dUpdated, err := svc.RegisterOrUpdate(ctx, RegisterInput{
		DeviceID:   "dev-abc-1",
		AppID:      "anr-001-wallpaper",
		Platform:   "android",
		OSVersion:  "15",
		AppVersion: "1.0.1",
	})
	if err != nil {
		t.Fatalf("unexpected error on update: %v", err)
	}

	if dUpdated.ID != 1 {
		t.Errorf("expected same ID 1 on upsert, got %d", dUpdated.ID)
	}
	if dUpdated.AppVersion != "1.0.1" {
		t.Errorf("expected updated version 1.0.1, got %s", dUpdated.AppVersion)
	}
}

func TestService_ValidationErrors(t *testing.T) {
	repo := NewMemoryRepository()
	svc := NewService(repo)
	ctx := context.Background()

	_, err := svc.RegisterOrUpdate(ctx, RegisterInput{
		DeviceID: "",
		AppID:    "anr-001-wallpaper",
		Platform: "android",
	})
	if !errors.Is(err, ErrInvalidDeviceID) {
		t.Errorf("expected ErrInvalidDeviceID, got %v", err)
	}

	_, err = svc.RegisterOrUpdate(ctx, RegisterInput{
		DeviceID: "dev-1",
		AppID:    "",
		Platform: "android",
	})
	if !errors.Is(err, ErrInvalidAppID) {
		t.Errorf("expected ErrInvalidAppID, got %v", err)
	}

	_, err = svc.RegisterOrUpdate(ctx, RegisterInput{
		DeviceID: "dev-1",
		AppID:    "anr-001-wallpaper",
		Platform: "windows",
	})
	if !errors.Is(err, ErrInvalidPlatform) {
		t.Errorf("expected ErrInvalidPlatform, got %v", err)
	}
}
