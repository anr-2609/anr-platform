package user

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/anr-2609/anr-platform/backend/internal/platform/auth"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrAccountInactive    = errors.New("account is inactive or suspended")
)

type Service struct {
	repo         Repository
	tokenService *auth.TokenService
}

func NewService(repo Repository, tokenService *auth.TokenService) *Service {
	return &Service{
		repo:         repo,
		tokenService: tokenService,
	}
}

func (s *Service) Authenticate(ctx context.Context, email, password string) (*User, string, string, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	if email == "" || password == "" {
		return nil, "", "", ErrInvalidCredentials
	}

	u, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil, "", "", ErrInvalidCredentials
		}
		return nil, "", "", fmt.Errorf("authenticating user: %w", err)
	}

	if u.Status != "active" {
		return nil, "", "", ErrAccountInactive
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return nil, "", "", ErrInvalidCredentials
	}

	accessToken, refreshToken, err := s.tokenService.GenerateUserTokens(u.ID, u.Email, u.Role)
	if err != nil {
		return nil, "", "", fmt.Errorf("generating user tokens: %w", err)
	}

	return u, accessToken, refreshToken, nil
}

func (s *Service) EnsureAdmin(ctx context.Context, defaultEmail, defaultPassword string) error {
	defaultEmail = strings.ToLower(strings.TrimSpace(defaultEmail))
	if defaultEmail == "" || defaultPassword == "" {
		return errors.New("default email and password required")
	}

	_, err := s.repo.FindByEmail(ctx, defaultEmail)
	if err == nil {
		return nil
	}

	if !errors.Is(err, ErrUserNotFound) {
		return fmt.Errorf("checking existing admin: %w", err)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(defaultPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hashing default admin password: %w", err)
	}

	adminUser := &User{
		Email:        defaultEmail,
		PasswordHash: string(hash),
		Role:         "admin",
		Status:       "active",
	}

	if err := s.repo.Create(ctx, adminUser); err != nil {
		return fmt.Errorf("creating default admin user: %w", err)
	}

	return nil
}

func (s *Service) FindByID(ctx context.Context, id int64) (*User, error) {
	return s.repo.FindByID(ctx, id)
}
