package auth

import (
	"github.com/alextilot/golang-htmx-chatapp/internal/config"
	"github.com/labstack/echo/v4"
	"log"
	"time"
)

// TokenRefresherMiddleware refreshes JWT tokens if the access token is about to expire
// and updates the UserContext in the request context.
func TokenRefresherMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		// Get access token from cookie
		accessCookie, err := ctx.Cookie(AccessTokenCookieName)
		if err != nil || accessCookie.Value == "" {
			// No token, proceed without user context
			return next(ctx)
		}

		// Parse claims from access token
		claims, err := ParseJWT(accessCookie.Value, []byte(config.Cfg.JwtSecretKey))
		if err != nil {
			// Invalid token, proceed without user context
			return next(ctx)
		}

		// We ensure that a new token is not issued until enough time has elapsed.
		// In this case, a new token will only be issued if the old token is within
		// 15 mins of expiry.
		if time.Until(time.Unix(claims.ExpiresAt, 0)) >= fifteenMinutes {
			return next(ctx)
		}

		// Gets the refresh token from the cookie.
		refreshCookie, err := ctx.Cookie(RefreshTokenCookieName)
		if err != nil || refreshCookie.Value == "" {
			return next(ctx)
		}

		_, err = ParseJWT(refreshCookie.Value, []byte(config.Cfg.JwtResfeshSecretKey))
		if err != nil {
			return next(ctx)
		}
		// Generate new tokens and update cookies
		_ = GenerateTokensAndSetCookies(claims, ctx)

		userCtx := FromClaims(claims)
		SetUserContext(ctx, userCtx)

		return next(ctx)
	}
}

func UserContextMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		accessCookie, err := c.Cookie(AccessTokenCookieName)
		if err != nil || accessCookie.Value == "" {
			SetUserContext(c, DefaultUserContext())
			return next(c)
		}

		claims, err := ParseJWT(accessCookie.Value, []byte(config.Cfg.JwtSecretKey))
		if err != nil {
			log.Printf("UserContextMiddleware: failed to parse JWT: %v", err)
			SetUserContext(c, DefaultUserContext())
			return next(c)
		}

		userCtx := FromClaims(claims)
		SetUserContext(c, userCtx)

		return next(c)
	}
}
