package user

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
)

type Repository interface {
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByID(ctx context.Context, id int64) (*User, error)
	Create(ctx context.Context, u *User) error
	Count(ctx context.Context) (int64, error)
}

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
	query := `
		SELECT id, email, password_hash, role, status, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	var u User
	err := r.pool.QueryRow(ctx, query, email).Scan(
		&u.ID,
		&u.Email,
		&u.PasswordHash,
		&u.Role,
		&u.Status,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("finding user by email: %w", err)
	}

	return &u, nil
}

func (r *PostgresRepository) FindByID(ctx context.Context, id int64) (*User, error) {
	query := `
		SELECT id, email, password_hash, role, status, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var u User
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&u.ID,
		&u.Email,
		&u.PasswordHash,
		&u.Role,
		&u.Status,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("finding user by id: %w", err)
	}

	return &u, nil
}

func (r *PostgresRepository) Create(ctx context.Context, u *User) error {
	now := time.Now().UTC()
	query := `
		INSERT INTO users (email, password_hash, role, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`

	err := r.pool.QueryRow(ctx, query,
		u.Email,
		u.PasswordHash,
		u.Role,
		u.Status,
		now,
		now,
	).Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt)

	if err != nil {
		return fmt.Errorf("creating user: %w", err)
	}

	return nil
}

func (r *PostgresRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM users`).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("counting users: %w", err)
	}
	return count, nil
}

type MemoryRepository struct {
	mu     sync.RWMutex
	users  map[int64]*User
	emails map[string]int64
	nextID int64
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		users:  make(map[int64]*User),
		emails: make(map[string]int64),
		nextID: 1,
	}
}

func (m *MemoryRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	id, ok := m.emails[email]
	if !ok {
		return nil, ErrUserNotFound
	}

	u := *m.users[id]
	return &u, nil
}

func (m *MemoryRepository) FindByID(ctx context.Context, id int64) (*User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	u, ok := m.users[id]
	if !ok {
		return nil, ErrUserNotFound
	}

	copyUser := *u
	return &copyUser, nil
}

func (m *MemoryRepository) Create(ctx context.Context, u *User) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.emails[u.Email]; ok {
		return ErrEmailAlreadyExists
	}

	now := time.Now().UTC()
	saved := *u
	saved.ID = m.nextID
	m.nextID++
	saved.CreatedAt = now
	saved.UpdatedAt = now

	m.users[saved.ID] = &saved
	m.emails[saved.Email] = saved.ID

	*u = saved
	return nil
}

func (m *MemoryRepository) Count(ctx context.Context) (int64, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return int64(len(m.users)), nil
}
