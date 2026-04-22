package auth

import (
	"github.com/alextilot/golang-htmx-chatapp/internal/auth/claims"
	"github.com/alextilot/golang-htmx-chatapp/internal/model"
	"github.com/alextilot/golang-htmx-chatapp/internal/usercontext"
	"github.com/labstack/echo/v4"
)

// CreateUserSession logs in a user by issuing tokens and attaching user context.
func CreateUserSession(user *model.User, c echo.Context) error {
	if user == nil {
		return nil
	}

	cl := claims.FromUser(user)

	// 1. Issue tokens (cookies + JWT)
	if err := IssueTokens(c, cl); err != nil {
		return err
	}

	// 2. Set request context directly (no re-derivation needed)
	usercontext.SetEcho(c, usercontext.FromClaims(cl))

	return nil
}
