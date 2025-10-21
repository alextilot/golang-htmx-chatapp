package auth

import (
	"github.com/alextilot/golang-htmx-chatapp/internal/auth/claims"
	"github.com/alextilot/golang-htmx-chatapp/internal/model"
	"github.com/alextilot/golang-htmx-chatapp/internal/usercontext"
	"github.com/labstack/echo/v4"
)

// SetUserAuthContext generates tokens, sets cookies, and updates UserContext
func SetUserAuthContext(user *model.User, c echo.Context) error {
	cl := claims.FromUser(user)
	if err := GenerateTokensAndSetCookies(cl, c); err != nil {
		return err
	}

	usercontext.Set(c, usercontext.FromClaims(cl))
	return nil
}
