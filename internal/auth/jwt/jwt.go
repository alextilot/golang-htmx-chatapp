package jwt

import (
	"github.com/alextilot/golang-htmx-chatapp/internal/auth/claims"
	"github.com/golang-jwt/jwt"
	"time"
)

// GenerateJWT generates a signed JWT string
func Generate(c *claims.Claims, ttl time.Duration, secret []byte) (string, time.Time, error) {
	exp := time.Now().Add(ttl)
	c.ExpiresAt = exp.Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	tknStr, err := token.SignedString(secret)
	return tknStr, exp, err
}

// ParseJWT parses a JWT string into Claims
func Parse(tokenStr string, secret []byte) (*claims.Claims, error) {
	c := &claims.Claims{}
	tkn, err := jwt.ParseWithClaims(tokenStr, c, func(token *jwt.Token) (any, error) {
		return secret, nil
	})
	if err != nil || !tkn.Valid {
		return nil, err
	}
	return c, nil
}
