package device

import "time"

type Device struct {
	ID           int64     `json:"id"`
	DeviceID     string    `json:"device_id"`
	AppID        string    `json:"app_id"`
	UserID       *int64    `json:"user_id,omitempty"`
	Platform     string    `json:"platform"`
	OSVersion    string    `json:"os_version"`
	AppVersion   string    `json:"app_version"`
	PushToken    string    `json:"push_token,omitempty"`
	LastActiveAt time.Time `json:"last_active_at"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
