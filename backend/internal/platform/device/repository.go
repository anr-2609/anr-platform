package device

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrDeviceNotFound = errors.New("device not found")
)

type Repository interface {
	Upsert(ctx context.Context, d *Device) (*Device, error)
	GetByDeviceAndApp(ctx context.Context, deviceID, appID string) (*Device, error)
	List(ctx context.Context, limit, offset int) ([]*Device, error)
	Count(ctx context.Context) (int64, error)
}

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) Upsert(ctx context.Context, d *Device) (*Device, error) {
	now := time.Now().UTC()
	query := `
		INSERT INTO devices (device_id, app_id, platform, os_version, app_version, push_token, last_active_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (device_id, app_id) DO UPDATE SET
			platform = EXCLUDED.platform,
			os_version = EXCLUDED.os_version,
			app_version = EXCLUDED.app_version,
			push_token = COALESCE(NULLIF(EXCLUDED.push_token, ''), devices.push_token),
			last_active_at = EXCLUDED.last_active_at,
			updated_at = EXCLUDED.updated_at
		RETURNING id, device_id, app_id, user_id, platform, os_version, app_version, push_token, last_active_at, created_at, updated_at
	`

	var saved Device
	err := r.pool.QueryRow(ctx, query,
		d.DeviceID,
		d.AppID,
		d.Platform,
		d.OSVersion,
		d.AppVersion,
		d.PushToken,
		now,
		now,
		now,
	).Scan(
		&saved.ID,
		&saved.DeviceID,
		&saved.AppID,
		&saved.UserID,
		&saved.Platform,
		&saved.OSVersion,
		&saved.AppVersion,
		&saved.PushToken,
		&saved.LastActiveAt,
		&saved.CreatedAt,
		&saved.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("upserting device: %w", err)
	}

	return &saved, nil
}

func (r *PostgresRepository) GetByDeviceAndApp(ctx context.Context, deviceID, appID string) (*Device, error) {
	query := `
		SELECT id, device_id, app_id, user_id, platform, os_version, app_version, push_token, last_active_at, created_at, updated_at
		FROM devices
		WHERE device_id = $1 AND app_id = $2
	`

	var d Device
	err := r.pool.QueryRow(ctx, query, deviceID, appID).Scan(
		&d.ID,
		&d.DeviceID,
		&d.AppID,
		&d.UserID,
		&d.Platform,
		&d.OSVersion,
		&d.AppVersion,
		&d.PushToken,
		&d.LastActiveAt,
		&d.CreatedAt,
		&d.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrDeviceNotFound
		}
		return nil, fmt.Errorf("getting device: %w", err)
	}

	return &d, nil
}

func (r *PostgresRepository) List(ctx context.Context, limit, offset int) ([]*Device, error) {
	if limit <= 0 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	query := `
		SELECT id, device_id, app_id, user_id, platform, os_version, app_version, push_token, last_active_at, created_at, updated_at
		FROM devices
		ORDER BY last_active_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("listing devices: %w", err)
	}
	defer rows.Close()

	var devices []*Device
	for rows.Next() {
		var d Device
		if err := rows.Scan(
			&d.ID,
			&d.DeviceID,
			&d.AppID,
			&d.UserID,
			&d.Platform,
			&d.OSVersion,
			&d.AppVersion,
			&d.PushToken,
			&d.LastActiveAt,
			&d.CreatedAt,
			&d.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scanning device: %w", err)
		}
		devices = append(devices, &d)
	}

	return devices, nil
}

func (r *PostgresRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM devices`).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("counting devices: %w", err)
	}
	return count, nil
}

type MemoryRepository struct {
	mu      sync.RWMutex
	devices map[string]*Device
	nextID  int64
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		devices: make(map[string]*Device),
		nextID:  1,
	}
}

func (m *MemoryRepository) key(deviceID, appID string) string {
	return deviceID + ":" + appID
}

func (m *MemoryRepository) Upsert(ctx context.Context, d *Device) (*Device, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	k := m.key(d.DeviceID, d.AppID)
	now := time.Now().UTC()

	existing, ok := m.devices[k]
	if ok {
		existing.Platform = d.Platform
		existing.OSVersion = d.OSVersion
		existing.AppVersion = d.AppVersion
		if d.PushToken != "" {
			existing.PushToken = d.PushToken
		}
		existing.LastActiveAt = now
		existing.UpdatedAt = now
		copyDevice := *existing
		return &copyDevice, nil
	}

	saved := *d
	saved.ID = m.nextID
	m.nextID++
	saved.CreatedAt = now
	saved.UpdatedAt = now
	saved.LastActiveAt = now

	m.devices[k] = &saved
	copyDevice := saved
	return &copyDevice, nil
}

func (m *MemoryRepository) GetByDeviceAndApp(ctx context.Context, deviceID, appID string) (*Device, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	k := m.key(deviceID, appID)
	d, ok := m.devices[k]
	if !ok {
		return nil, ErrDeviceNotFound
	}

	copyDevice := *d
	return &copyDevice, nil
}

func (m *MemoryRepository) List(ctx context.Context, limit, offset int) ([]*Device, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var list []*Device
	for _, d := range m.devices {
		copyDevice := *d
		list = append(list, &copyDevice)
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].LastActiveAt.After(list[j].LastActiveAt)
	})

	if offset >= len(list) {
		return []*Device{}, nil
	}

	end := offset + limit
	if limit <= 0 || end > len(list) {
		end = len(list)
	}

	return list[offset:end], nil
}

func (m *MemoryRepository) Count(ctx context.Context) (int64, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return int64(len(m.devices)), nil
}
