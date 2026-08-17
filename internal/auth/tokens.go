package auth

import (
	"time"
)

const (
	AccessTokenTTL  = 1 * time.Hour
	RefreshTokenTTL = 24 * time.Hour
)

func (s *Service) GenerateAccessToken(c *Claims) (string, time.Time, error) {
	return generateJWT(c, AccessTokenTTL, s.cfg.JWTSecretKey)
}

func (s *Service) GenerateRefreshToken(c *Claims) (string, time.Time, error) {
	return generateJWT(c, RefreshTokenTTL, s.cfg.JWTRefreshSecretKey)
}
