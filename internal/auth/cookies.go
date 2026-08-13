package auth

import (
	"net/http"
	"time"

	"github.com/alextilot/golang-htmx-chatapp/internal/config"
	"github.com/labstack/echo/v5"
)

const (
	AccessTokenCookieName  = "access-token"
	RefreshTokenCookieName = "refresh-token"
)

// setCookie applies consistent secure cookie configuration.
func setCookie(c *echo.Context, name string, value string, expires time.Time) {
	c.SetCookie(&http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		Expires:  expires,
		HttpOnly: true,

		// Only send cookies over HTTPS in production
		Secure: config.Cfg.IsProduction(),

		// Good default for auth flows (login redirects, same-site requests)
		SameSite: http.SameSiteLaxMode,
	})
}

// SetAccessToken sets the access token cookie.
func SetAccessToken(c *echo.Context, token string, expires time.Time) {
	setCookie(c, AccessTokenCookieName, token, expires)
}

// SetRefreshToken sets the refresh token cookie.
func SetRefreshToken(c *echo.Context, token string, expires time.Time) {
	setCookie(c, RefreshTokenCookieName, token, expires)
}

// Clear removes both authentication cookies.
func Clear(c *echo.Context) {
	clear := func(name string) {
		c.SetCookie(&http.Cookie{
			Name:     name,
			Value:    "",
			Path:     "/",
			Expires:  time.Unix(0, 0),
			HttpOnly: true,
			Secure:   config.Cfg.IsProduction(),
			SameSite: http.SameSiteLaxMode,
		})
	}

	clear(AccessTokenCookieName)
	clear(RefreshTokenCookieName)
}

// GetCookieValue safely retrieves a cookie value.
func GetCookieValue(c *echo.Context, name string) (string, error) {
	cookie, err := c.Cookie(name)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}
