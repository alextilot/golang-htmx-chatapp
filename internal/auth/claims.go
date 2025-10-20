package auth

import (
	"github.com/alextilot/golang-htmx-chatapp/internal/model"
	"github.com/golang-jwt/jwt"
	"time"
)

type Claims struct {
	UserID   string `json:"userID"`
	Username string `json:"username"`
	Email    string `json:"email"`
	// Theme       string   `json:"theme,omitempty"`
	// Permissions []string `json:"permissions,omitempty"`
	jwt.StandardClaims
}

func GenerateJWT(claims *Claims, ttl time.Duration, secret []byte) (string, time.Time, error) {
	exp := time.Now().Add(ttl)
	claims.ExpiresAt = exp.Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tknStr, err := token.SignedString(secret)
	return tknStr, exp, err
}

func ParseJWT(tokenStr string, secret []byte) (*Claims, error) {
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (any, error) {
		return secret, nil
	})
	if err != nil || !tkn.Valid {
		return nil, err
	}
	return claims, nil
}

func ClaimsFromUser(user *model.User) *Claims {
	return &Claims{
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
	}
}
