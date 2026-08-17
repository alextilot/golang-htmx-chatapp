package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// generateJWT signs c into a JWT valid for ttl, returning the token and its
// expiry.
func generateJWT(c *Claims, ttl time.Duration, secret []byte) (string, time.Time, error) {
	now := time.Now()
	exp := now.Add(ttl)

	claimsCopy := *c
	claimsCopy.ExpiresAt = jwt.NewNumericDate(exp)
	claimsCopy.IssuedAt = jwt.NewNumericDate(now)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claimsCopy)
	tknStr, err := token.SignedString(secret)

	return tknStr, exp, err
}

// parseJWT parses and verifies a JWT string, returning its Claims.
func parseJWT(tokenStr string, secret []byte) (*Claims, error) {
	c := &Claims{}

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
