package auth

import (
	"github.com/alextilot/golang-htmx-chatapp/internal/auth/claims"
	"github.com/alextilot/golang-htmx-chatapp/internal/model"
	"github.com/labstack/echo/v5"
)

// CreateUserSession logs in a user by issuing tokens and attaching the
// authenticated principal to the request context.
func (s *Service) CreateUserSession(user *model.User, c *echo.Context) error {
	if user == nil {
		return nil
	}

	cl := claims.FromUser(user)

	// 1. Issue tokens (cookies + JWT)
	if err := s.IssueTokens(c, cl); err != nil {
		return err
	}

	// 2. Set request context directly (no re-derivation needed)
	SetEchoPrincipal(c, PrincipalFromClaims(cl))

	return nil
}
