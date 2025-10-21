package usercontext

import (
	"context"
	"github.com/alextilot/golang-htmx-chatapp/internal/auth/claims"
	"github.com/labstack/echo/v4"
)

// UserContext holds all per-request user state.
// It can include both authentication info and user-specific preferences.
type UserContext struct {
	ID         string
	Username   string
	Email      string
	IsLoggedIn bool
	// Theme       string
	// Locale      string
	// Permissions []string
}

type contextKey struct{}

// key is an unexported zero-value type to avoid context key collisions.
var key contextKey

func Default() *UserContext {
	return &UserContext{IsLoggedIn: false}
}

// FromClaims creates a new UserContext with default values or from claims
func FromClaims(claims *claims.Claims) *UserContext {
	if claims == nil {
		return Default()
	}

	return &UserContext{
		ID:         claims.UserID,
		Username:   claims.Username,
		Email:      claims.Email,
		IsLoggedIn: true,
	}
}

// With returns a new context containing the given UserContext.
func With(ctx context.Context, u *UserContext) context.Context {
	return context.WithValue(ctx, key, u)
}

// From extracts the UserContext from a context.Context.
// Returns Default() if none is found.
func From(ctx context.Context) *UserContext {
	if v, ok := ctx.Value(key).(*UserContext); ok && v != nil {
		return v
	}
	return Default()
}

// Get retrieves the UserContext from an Echo context.
func Get(ctx context.Context) *UserContext {
	if v, ok := ctx.Value(key).(*UserContext); ok {
		return v
	}
	return FromClaims(nil)
}

// Set attaches a UserContext to an Echo context.
func Set(c echo.Context, u *UserContext) {
	c.SetRequest(c.Request().WithContext(With(c.Request().Context(), u)))
}
