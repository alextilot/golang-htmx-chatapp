package auth

import (
	"github.com/alextilot/golang-htmx-chatapp/internal/auth/claims"
	"github.com/alextilot/golang-htmx-chatapp/internal/auth/jwt"
	"github.com/alextilot/golang-htmx-chatapp/internal/config"
	"time"
)

const (
	AccessTokenTTL  = 1 * time.Hour
	RefreshTokenTTL = 24 * time.Hour
)

func GenerateAccessToken(c *claims.Claims) (string, time.Time, error) {
	return jwt.Generate(c, AccessTokenTTL, []byte(config.Cfg.JwtSecretKey))
}

func GenerateRefreshToken(c *claims.Claims) (string, time.Time, error) {
	return jwt.Generate(c, RefreshTokenTTL, []byte(config.Cfg.JwtRefeshSecretKey))
}
