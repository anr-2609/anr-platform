package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("invalid or expired token")
	ErrInvalidType  = errors.New("invalid token type")
)

type DeviceClaims struct {
	DeviceID string `json:"device_id"`
	AppID    string `json:"app_id"`
	Type     string `json:"type"`
	jwt.RegisteredClaims
}

type TokenService struct {
	secret     []byte
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func NewTokenService(secret string, accessTTL, refreshTTL time.Duration) *TokenService {
	return &TokenService{
		secret:     []byte(secret),
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
	}
}

func (s *TokenService) GenerateTokens(deviceID, appID string) (string, string, error) {
	now := time.Now().UTC()

	accessClaims := DeviceClaims{
		DeviceID: deviceID,
		AppID:    appID,
		Type:     "GUEST_ACCESS",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   deviceID,
			Issuer:    "anr-platform",
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.accessTTL)),
		},
	}

	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString(s.secret)
	if err != nil {
		return "", "", fmt.Errorf("signing access token: %w", err)
	}

	refreshClaims := DeviceClaims{
		DeviceID: deviceID,
		AppID:    appID,
		Type:     "GUEST_REFRESH",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   deviceID,
			Issuer:    "anr-platform",
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.refreshTTL)),
		},
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString(s.secret)
	if err != nil {
		return "", "", fmt.Errorf("signing refresh token: %w", err)
	}

	return accessToken, refreshToken, nil
}

func (s *TokenService) ValidateAccessToken(tokenStr string) (*DeviceClaims, error) {
	claims, err := s.parseClaims(tokenStr)
	if err != nil {
		return nil, err
	}

	if claims.Type != "GUEST_ACCESS" {
		return nil, ErrInvalidType
	}

	return claims, nil
}

func (s *TokenService) ValidateRefreshToken(tokenStr string) (*DeviceClaims, error) {
	claims, err := s.parseClaims(tokenStr)
	if err != nil {
		return nil, err
	}

	if claims.Type != "GUEST_REFRESH" {
		return nil, ErrInvalidType
	}

	return claims, nil
}

func (s *TokenService) parseClaims(tokenStr string) (*DeviceClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &DeviceClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.secret, nil
	})

	if err != nil || !token.Valid {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*DeviceClaims)
	if !ok {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
