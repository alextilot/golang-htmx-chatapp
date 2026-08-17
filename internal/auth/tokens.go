package auth

import (
	"time"

	"github.com/alextilot/golang-htmx-chatapp/internal/auth/claims"
	"github.com/alextilot/golang-htmx-chatapp/internal/auth/jwt"
)

const (
	AccessTokenTTL  = 1 * time.Hour
	RefreshTokenTTL = 24 * time.Hour
)

func (s *Service) GenerateAccessToken(c *claims.Claims) (string, time.Time, error) {
	return jwt.Generate(c, AccessTokenTTL, s.cfg.JWTSecretKey)
}

func (s *Service) GenerateRefreshToken(c *claims.Claims) (string, time.Time, error) {
	return jwt.Generate(c, RefreshTokenTTL, s.cfg.JWTRefreshSecretKey)
}
