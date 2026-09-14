package user

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/anr-2609/anr-platform/backend/internal/platform/auth"
)

func TestService_Authenticate_And_EnsureAdmin(t *testing.T) {
	ctx := context.Background()
	repo := NewMemoryRepository()
	tokenService := auth.NewTokenService("test_secret_for_unit_tests_32_chars_long", 15*time.Minute, 7*24*time.Hour)
	svc := NewService(repo, tokenService)

	err := svc.EnsureAdmin(ctx, "admin@anr-studio.com", "SecureAdminPass123!")
	require.NoError(t, err)

	err = svc.EnsureAdmin(ctx, "admin@anr-studio.com", "SecureAdminPass123!")
	require.NoError(t, err)

	count, err := repo.Count(ctx)
	require.NoError(t, err)
	assert.Equal(t, int64(1), count)

	u, accessToken, refreshToken, err := svc.Authenticate(ctx, "admin@anr-studio.com", "SecureAdminPass123!")
	require.NoError(t, err)
	assert.NotNil(t, u)
	assert.Equal(t, "admin@anr-studio.com", u.Email)
	assert.Equal(t, "admin", u.Role)
	assert.NotEmpty(t, accessToken)
	assert.NotEmpty(t, refreshToken)

	claims, err := tokenService.ValidateUserAccessToken(accessToken)
	require.NoError(t, err)
	assert.Equal(t, u.ID, claims.UserID)
	assert.Equal(t, "admin@anr-studio.com", claims.Email)
	assert.Equal(t, "admin", claims.Role)

	_, _, _, err = svc.Authenticate(ctx, "admin@anr-studio.com", "WrongPassword")
	assert.ErrorIs(t, err, ErrInvalidCredentials)

	_, _, _, err = svc.Authenticate(ctx, "nonexistent@anr-studio.com", "SomePassword")
	assert.ErrorIs(t, err, ErrInvalidCredentials)
}
