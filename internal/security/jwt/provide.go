package jwt

import (
	"time"

	"github.com/Koshsky/erp-backend/internal/config"
)

// ProvideJWTService builds the JWT service.
//
// refreshKey (JWT_REFRESH_KEY) больше не используется: refresh-токены opaque и
// хранятся в БД (AD-06), ключ сохранён в конфиге/окружении для совместимости.
func ProvideJWTService(cfg config.JWTConfig) *Service {
	return &Service{
		secretKey:     []byte(cfg.SecretKey),
		accessExpiry:  time.Duration(cfg.AccessExpiry),
		refreshExpiry: time.Duration(cfg.RefreshExpiry),
		issuer:        cfg.Issuer,
	}
}
