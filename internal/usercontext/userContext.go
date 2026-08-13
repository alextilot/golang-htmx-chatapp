package usercontext

import (
	"context"

	"github.com/alextilotecho/v5-htmx-chatapp/internal/auth/claims"
	"github.com/labstack/echo/v5"
)

// UserContext holds all per-request user state.
// It can include both authentication info and user-specific preferences.
type UserContext struct {
	ID            string
	Username      string
	Email         string
	Authenticated bool
	// Theme       string
	// Locale      string
	// Permissions []string
}

type contextKey struct{}

// key is an unexported zero-value type to avoid context key collisions.
var key contextKey

var anonymous = &UserContext{Authenticated: false}

func Default() *UserContext {
	return anonymous
}

func (u *UserContext) IsAnonymous() bool {
	return !u.Authenticated
}

// FromClaims creates a new UserContext with default values or from claims
func FromClaims(claims *claims.Claims) *UserContext {
	if claims == nil {
		return Default()
	}

	return &UserContext{
		ID:       claims.UserID,
		Username: claims.Username,
		Email:    claims.Email,
		// malformed token could mark someone as authenticated.
		Authenticated: claims.UserID != "",
	}
}

func SetContext(ctx context.Context, u *UserContext) context.Context {
	return context.WithValue(ctx, key, u)
}

// FromContext extracts the UserContext from a context.Context.
// Returns Default() if none is found.
func FromContext(ctx context.Context) *UserContext {
	if v, ok := ctx.Value(key).(*UserContext); ok && v != nil {
		return v
	}
	return Default()
}

// SetEcho attaches a UserContext to an Echo context.
func SetEcho(c *echo.Context, u *UserContext) {
	req := c.Request()
	ctx := SetContext(req.Context(), u)
	c.SetRequest(req.WithContext(ctx))
}

// FromEcho retrieves the UserContext from an Echo context.
func FromEcho(c *echo.Context) *UserContext {
	return FromContext(c.Request().Context())
}
