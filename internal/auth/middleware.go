package auth

import (
	"log"
	"time"

	"github.com/alextilot/golang-htmx-chatapp/internal/auth/claims"
	"github.com/alextilot/golang-htmx-chatapp/internal/auth/jwt"
	"github.com/alextilot/golang-htmx-chatapp/internal/config"
	"github.com/alextilot/golang-htmx-chatapp/internal/usercontext"
	"github.com/labstack/echo/v5"
)

const fifteenMinutes = 15 * time.Minute

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		uc := usercontext.Default()

		accessToken, err := GetCookieValue(c, AccessTokenCookieName)
		if err != nil || accessToken == "" {
			usercontext.SetEcho(c, uc)
			return next(c)
		}

		parsedClaims, err := jwt.Parse(accessToken, []byte(config.Cfg.JwtSecretKey))
		if err != nil {
			log.Printf("AuthMiddleware: invalid access token: %v", err)
			usercontext.SetEcho(c, uc)
			return next(c)
		}

		uc = usercontext.FromClaims(parsedClaims)

		tryRefresh(c, parsedClaims)

		usercontext.SetEcho(c, uc)
		return next(c)
	}
}

func IssueTokens(c *echo.Context, userClaims *claims.Claims) error {
	accessToken, accessExp, err := GenerateAccessToken(userClaims)
	if err != nil {
		return err
	}

	SetAccessToken(c, accessToken, accessExp)

	refreshToken, refreshExp, err := GenerateRefreshToken(userClaims)
	if err != nil {
		return err
	}

	SetRefreshToken(c, refreshToken, refreshExp)

	return nil
}

func tryRefresh(c *echo.Context, cClaims *claims.Claims) {
	if !isExpiringSoon(cClaims.ExpiresAt.Unix()) {
		return
	}

	refreshToken, err := GetCookieValue(c, RefreshTokenCookieName)
	if err != nil || refreshToken == "" {
		return
	}

	_, err = jwt.Parse(refreshToken, []byte(config.Cfg.JwtRefeshSecretKey))
	if err != nil {
		return
	}

	if err := IssueTokens(c, cClaims); err != nil {
		log.Printf("AuthMiddleware: failed to refresh tokens: %v", err)
	}
}

func isExpiringSoon(exp int64) bool {
	return time.Until(time.Unix(exp, 0)) < fifteenMinutes
}

func RequireLogin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		if !usercontext.FromEcho(c).Authenticated {
			return c.JSON(401, map[string]string{
				"error": "unauthorized: please log in",
			})
		}
		return next(c)
	}
}
