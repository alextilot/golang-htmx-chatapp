package auth

import (
	"context"

	"github.com/labstack/echo/v5"
)

// Principal holds the authenticated request identity derived from a JWT.
//
// Principal is a snapshot of claims at token-issue time. It is not a place
// for arbitrary user preferences or domain data — those belong to the
// domain user (loaded from the service/repository layer when needed), not
// to the authentication identity carried on every request.
type Principal struct {
	ID            string
	Username      string
	Email         string
	Authenticated bool
}

type principalKey struct{}

var key principalKey

// anonymousPrincipal is the zero-value, unauthenticated identity attached to
// every request until/unless AuthMiddleware replaces it.
var anonymousPrincipal = &Principal{Authenticated: false}

// AnonymousPrincipal returns the shared unauthenticated principal.
func AnonymousPrincipal() *Principal {
	return anonymousPrincipal
}

// IsAnonymous reports whether this principal is unauthenticated.
func (p *Principal) IsAnonymous() bool {
	return !p.Authenticated
}

// PrincipalFromClaims builds a Principal from parsed JWT claims.
func PrincipalFromClaims(c *Claims) *Principal {
	if c == nil {
		return AnonymousPrincipal()
	}

	return &Principal{
		ID:       c.UserID,
		Username: c.Username,
		Email:    c.Email,
		// A malformed/empty-subject token should never be treated as authenticated.
		Authenticated: c.UserID != "",
	}
}

// WithPrincipal returns a new context carrying the given Principal.
func WithPrincipal(ctx context.Context, p *Principal) context.Context {
	return context.WithValue(ctx, key, p)
}

// PrincipalFromContext extracts the Principal from ctx, or the anonymous
// principal if none was attached.
func PrincipalFromContext(ctx context.Context) *Principal {
	if v, ok := ctx.Value(key).(*Principal); ok && v != nil {
		return v
	}
	return AnonymousPrincipal()
}

// SetEchoPrincipal attaches a Principal to an Echo request context.
func SetEchoPrincipal(c *echo.Context, p *Principal) {
	req := c.Request()
	c.SetRequest(req.WithContext(WithPrincipal(req.Context(), p)))
}

// PrincipalFromEcho retrieves the Principal from an Echo request context.
func PrincipalFromEcho(c *echo.Context) *Principal {
	return PrincipalFromContext(c.Request().Context())
}
