package jwt

import (
	"github.com/alextilot/golang-htmx-chatapp/internal/auth/claims"
	"github.com/golang-jwt/jwt"
	"time"
)

func Generate(c *claims.Claims, ttl time.Duration, secret []byte) (string, time.Time, error) {
	now := time.Now()
	exp := now.Add(ttl)

	claimsCopy := *c
	claimsCopy.ExpiresAt = exp.Unix()
	claimsCopy.IssuedAt = now.Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claimsCopy)
	tknStr, err := token.SignedString(secret)

	return tknStr, exp, err
}

// ParseJWT parses a JWT string into Claims
func Parse(tokenStr string, secret []byte) (*claims.Claims, error) {
	c := &claims.Claims{}

	tkn, err := jwt.ParseWithClaims(tokenStr, c, func(token *jwt.Token) (any, error) {
		// Enforce HMAC signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return secret, nil
	})

	if err != nil || !tkn.Valid {
		return nil, err
	}

	return c, nil
}
