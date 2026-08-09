package jwt

import (
	"time"

	"github.com/Koshsky/erp-backend/internal/config"
)

// ProvideJWTService builds the JWT service.
func ProvideJWTService(cfg config.JWTConfig) *Service {
	return &Service{
		secretKey:     []byte(cfg.SecretKey),
		refreshKey:    []byte(cfg.RefreshKey),
		accessExpiry:  time.Duration(cfg.AccessExpiry),
		refreshExpiry: time.Duration(cfg.RefreshExpiry),
		issuer:        cfg.Issuer,
	}
}
