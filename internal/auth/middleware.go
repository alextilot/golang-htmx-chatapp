package auth

import (
	"github.com/alextilot/golang-htmx-chatapp/internal/auth/jwt"
	"github.com/alextilot/golang-htmx-chatapp/internal/config"
	"github.com/alextilot/golang-htmx-chatapp/internal/usercontext"
	"github.com/labstack/echo/v4"
	"log"
	"time"
)

const fifteenMinutes = 15 * time.Minute

// TokenRefresherMiddleware refreshes JWT tokens if the access token is about to expire
// and updates the UserContext in the request context.
func TokenRefresherMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		// Get access token from cookie
		accessCookie, err := ctx.Cookie(AccessTokenCookieName)
		if err != nil || accessCookie.Value == "" {
			// No access token — continue anonymously
			return next(ctx)
		}

		// Parse access token claims
		claims, err := jwt.Parse(accessCookie.Value, []byte(config.Cfg.JwtSecretKey))
		if err != nil {
			log.Printf("TokenRefresherMiddleware: invalid access token: %v", err)
			return next(ctx)
		}

		// Refresh token only if the access token expires within 15 minutes
		if time.Until(time.Unix(claims.ExpiresAt, 0)) >= fifteenMinutes {
			return next(ctx)
		}

		// Gets refresh token from cookie
		refreshCookie, err := ctx.Cookie(RefreshTokenCookieName)
		if err != nil || refreshCookie.Value == "" {
			return next(ctx)
		}

		// Parse refresh token claims
		_, err = jwt.Parse(refreshCookie.Value, []byte(config.Cfg.JwtRefeshSecretKey))
		if err != nil {
			return next(ctx)
		}

		if err := GenerateAndSet(ctx, claims); err != nil {
			log.Printf("TokenRefresherMiddleware: failed to refresh tokens: %v", err)
		}

		// Always attach user context (even if not refreshed)
		usercontext.Set(ctx, usercontext.FromClaims(claims))

		return next(ctx)
	}
}

// UserContextMiddleware attaches a UserContext to the request context based on JWT claims.
// If no valid access token is found, a default anonymous context is used.
func UserContextMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Get access token from cookie
		accessCookie, err := c.Cookie(AccessTokenCookieName)
		if err != nil || accessCookie.Value == "" {
			// No token → set default anonymous context
			usercontext.Set(c, usercontext.Default())
			return next(c)
		}

		claims, err := jwt.Parse(accessCookie.Value, []byte(config.Cfg.JwtSecretKey))
		if err != nil {
			// Invalid token → log and set default context
			log.Printf("UserContextMiddleware: failed to parse JWT: %v", err)
			usercontext.Set(c, usercontext.Default())
			return next(c)
		}

		usercontext.Set(c, usercontext.FromClaims(claims))

		return next(c)
	}
}
