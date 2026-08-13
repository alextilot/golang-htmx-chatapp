package claims

import (
	"github.com/alextilot/golang-htmx-chatapp/internal/model"
	"github.com/golang-jwt/jwt/v5"
)

// Claims represents the JWT payload
type Claims struct {
	UserID   string `json:"userID"`
	Username string `json:"username"`
	Email    string `json:"email"`
	// Theme       string   `json:"theme,omitempty"`
	// Permissions []string `json:"permissions,omitempty"`
	jwt.RegisteredClaims
}

// FromUser creates Claims from a User model
func FromUser(user *model.User) *Claims {
	if user == nil {
		return &Claims{}
	}

	return &Claims{
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
	}
}
