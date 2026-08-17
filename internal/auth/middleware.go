package auth

import (
	"log"
	"time"

	"github.com/labstack/echo/v5"
)

const fifteenMinutes = 15 * time.Minute

// AuthMiddleware attaches an auth.Principal (authenticated or anonymous) to
// every request based on the access token cookie, refreshing it if it's
// expiring soon.
func (s *Service) AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		p := AnonymousPrincipal()

		accessToken, err := GetCookieValue(c, AccessTokenCookieName)
		if err != nil || accessToken == "" {
			SetEchoPrincipal(c, p)
			return next(c)
		}

		parsedClaims, err := parseJWT(accessToken, s.cfg.JWTSecretKey)
		if err != nil {
			log.Printf("AuthMiddleware: invalid access token: %v", err)
			SetEchoPrincipal(c, p)
			return next(c)
		}

		p = PrincipalFromClaims(parsedClaims)

		s.tryRefresh(c, parsedClaims)

		SetEchoPrincipal(c, p)
		return next(c)
	}
}

// IssueTokens generates and sets both the access and refresh token cookies
// for the given claims.
func (s *Service) IssueTokens(c *echo.Context, userClaims *Claims) error {
	accessToken, accessExp, err := s.GenerateAccessToken(userClaims)
	if err != nil {
		return err
	}

	s.SetAccessToken(c, accessToken, accessExp)

	refreshToken, refreshExp, err := s.GenerateRefreshToken(userClaims)
	if err != nil {
		return err
	}

	s.SetRefreshToken(c, refreshToken, refreshExp)

	return nil
}

func (s *Service) tryRefresh(c *echo.Context, cClaims *Claims) {
	if !isExpiringSoon(cClaims.ExpiresAt.Unix()) {
		return
	}

	refreshToken, err := GetCookieValue(c, RefreshTokenCookieName)
	if err != nil || refreshToken == "" {
		return
	}

	if _, err := parseJWT(refreshToken, s.cfg.JWTRefreshSecretKey); err != nil {
		return
	}

	if err := s.IssueTokens(c, cClaims); err != nil {
		log.Printf("AuthMiddleware: failed to refresh tokens: %v", err)
	}
}

func isExpiringSoon(exp int64) bool {
	return time.Until(time.Unix(exp, 0)) < fifteenMinutes
}

// RequireLogin rejects unauthenticated requests. It only inspects the
// Principal already attached by AuthMiddleware, so it needs no
// configuration and stays a free function rather than a Service method.
func RequireLogin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		if !PrincipalFromEcho(c).Authenticated {
			return c.JSON(401, map[string]string{
				"error": "unauthorized: please log in",
			})
		}
		return next(c)
	}
}
