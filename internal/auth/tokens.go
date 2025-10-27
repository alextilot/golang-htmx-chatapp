package auth

import (
	"github.com/alextilot/golang-htmx-chatapp/internal/auth/claims"
	"github.com/alextilot/golang-htmx-chatapp/internal/auth/jwt"
	"github.com/alextilot/golang-htmx-chatapp/internal/config"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

const (
	AccessTokenCookieName  = "access-token"
	RefreshTokenCookieName = "refresh-token"

	accessTokenTTL  = 1 * time.Hour
	refreshTokenTTL = 24 * time.Hour
)

// setCookie is an internal helper that applies consistent cookie configuration.
func setCookie(c echo.Context, name string, value string, expires time.Time) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		Expires:  expires,
		HttpOnly: true,
		Secure:   config.Cfg.IsProduction(), // only send over HTTPS in prod
		SameSite: http.SameSiteLaxMode,      // prevent CSRF while allowing login redirects
	}
	c.SetCookie(cookie)
}

// Clear removes both auth cookies (access + refresh).
func Clear(ctx echo.Context) {
	clearCookie := func(name string) {
		ctx.SetCookie(&http.Cookie{
			Name:     name,
			Value:    "",
			Path:     "/",
			Expires:  time.Unix(0, 0),
			HttpOnly: true,
			Secure:   config.Cfg.IsProduction(),
			SameSite: http.SameSiteLaxMode,
		})
	}

	clearCookie(AccessTokenCookieName)
	clearCookie(RefreshTokenCookieName)
}

// GenerateAndSet creates JWT tokens (access + refresh) and sets them as cookies.
func GenerateAndSet(ctx echo.Context, userClaims *claims.Claims) error {
	accessToken, accessExp, err := jwt.Generate(userClaims, accessTokenTTL, []byte(config.Cfg.JwtSecretKey))
	if err != nil {
		return err
	}
	setCookie(ctx, AccessTokenCookieName, accessToken, accessExp)

	refreshToken, refreshExp, err := jwt.Generate(userClaims, refreshTokenTTL, []byte(config.Cfg.JwtRefeshSecretKey))
	if err != nil {
		return err
	}
	setCookie(ctx, RefreshTokenCookieName, refreshToken, refreshExp)

	return nil
}

// Get retrieves a specific cookie value safely.
func Get(ctx echo.Context, name string) (string, error) {
	cookie, err := ctx.Cookie(name)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}
