package auth

import (
	"github.com/alextilot/golang-htmx-chatapp/internal/config"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

const (
	AccessTokenCookieName  = "access-token"
	RefreshTokenCookieName = "refresh-token"

	oneHour         = 1 * time.Hour
	twentyFourHours = 24 * time.Hour
	fifteenMinutes  = 15 * time.Minute
)

func setTokenCookie(name string, token string, expiration time.Time, ctx echo.Context) {
	c := &http.Cookie{
		Name:     name,
		Value:    token,
		Expires:  expiration,
		Path:     "/",
		HttpOnly: true,
		// Secure:   config.Cfg.IsProduction,
		// SameSite: http.SameSiteLaxMode,
	}
	ctx.SetCookie(c)
}

func GenerateTokensAndSetCookies(claims *Claims, ctx echo.Context) error {
	accessToken, exp, err := GenerateJWT(claims, oneHour, []byte(config.Cfg.JwtSecretKey))
	if err != nil {
		return err
	}
	setTokenCookie(AccessTokenCookieName, accessToken, exp, ctx)

	refreshToken, exp, err := GenerateJWT(claims, twentyFourHours, []byte(config.Cfg.JwtResfeshSecretKey))
	if err != nil {
		return err
	}
	setTokenCookie(RefreshTokenCookieName, refreshToken, exp, ctx)

	return nil
}
