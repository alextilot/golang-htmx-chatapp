package auth

import (
	"context"
	"github.com/labstack/echo/v4"
)

type UserContext struct {
	ID         string
	Username   string
	Email      string
	IsLoggedIn bool
	// Theme       string
	// Permissions []string
}

type contextKey string

const userContextKey contextKey = "userContext"

func DefaultUserContext() *UserContext {
	return &UserContext{IsLoggedIn: false}
}

// FromClaims creates a new UserContext with default values or from claims
func FromClaims(claims *Claims) *UserContext {
	if claims == nil {
		return DefaultUserContext()
	}

	return &UserContext{
		ID:         claims.UserID,
		Username:   claims.Username,
		Email:      claims.Email,
		IsLoggedIn: true,
	}
}

func WithUserContext(ctx context.Context, u *UserContext) context.Context {
	return context.WithValue(ctx, userContextKey, u)
}

func GetUserContext(ctx context.Context) *UserContext {
	if v, ok := ctx.Value(userContextKey).(*UserContext); ok {
		return v
	}
	return FromClaims(nil)
}

// SetUserContext is a helper to set UserContext on an Echo context
func SetUserContext(c echo.Context, u *UserContext) {
	c.SetRequest(c.Request().WithContext(WithUserContext(c.Request().Context(), u)))
}
