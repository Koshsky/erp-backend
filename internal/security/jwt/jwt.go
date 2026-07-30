package jwt

import (
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/Koshsky/erp-backend/internal/config"
)

type Service struct {
	secretKey     []byte
	refreshKey    []byte
	accessExpiry  time.Duration
	refreshExpiry time.Duration
	issuer        string
}

func NewJWTService(config config.JWTConfig) *Service {
	return &Service{
		secretKey:     []byte(config.SecretKey),
		refreshKey:    []byte(config.RefreshKey),
		accessExpiry:  config.AccessExpiry,
		refreshExpiry: config.RefreshExpiry,
		issuer:        config.Issuer,
	}
}

// GenerateAccessToken generates a new JWT access token for the given user ID and role.
func (s *Service) GenerateAccessToken(userID int64, role, email string) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID: userID,
		Role:   role,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(s.accessExpiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    s.issuer,
			Subject:   strconv.FormatInt(userID, 10),
			ID:        generateTokenID(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secretKey)
}

// GenerateRefreshToken generates a refresh token with the given user ID.
func (s *Service) GenerateRefreshToken(userID int64) (string, error) {
	now := time.Now()
	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(now.Add(s.refreshExpiry)),
		IssuedAt:  jwt.NewNumericDate(now),
		NotBefore: jwt.NewNumericDate(now),
		Issuer:    s.issuer,
		Subject:   strconv.FormatInt(userID, 10),
		ID:        generateTokenID(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.refreshKey)
}

// GenerateTokenPair generates both access and refresh tokens
func (s *Service) GenerateTokenPair(userID int64, role, email string) (*TokenPair, error) {
	accessToken, err := s.GenerateAccessToken(userID, role, email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.GenerateRefreshToken(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int(s.accessExpiry.Seconds()),
	}, nil
}

// ValidateAccessToken проверяет и парсит access токен
func (s *Service) ValidateAccessToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.secretKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("invalid access token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid access token claims")
	}

	return claims, nil
}

// ValidateRefreshToken проверяет и парсит refresh токен
func (s *Service) ValidateRefreshToken(tokenString string) (*jwt.RegisteredClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.refreshKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid refresh token claims")
	}

	return claims, nil
}

// RefreshAccessToken создаёт новый access token по refresh токену
func (s *Service) RefreshAccessToken(refreshTokenString string) (*TokenPair, error) {
	claims, err := s.ValidateRefreshToken(refreshTokenString)
	if err != nil {
		return nil, err
	}

	// Извлекаем userID из Subject
	var userID int64
	if _, err := fmt.Sscanf(claims.Subject, "%d", &userID); err != nil {
		return nil, fmt.Errorf("invalid user id in token")
	}

	// Генерируем только новый access token
	accessToken, err := s.GenerateAccessToken(userID, "", "")
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenString,
		TokenType:    "Bearer",
		ExpiresIn:    int(s.accessExpiry.Seconds()),
	}, nil
}

func generateTokenID() string {
	return strconv.FormatInt(time.Now().UnixNano(), 10)
}
