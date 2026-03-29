package model

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Base

	Username string `gorm:"unique;not null"` // Unique and non-nullable username
	Email    string `gorm:"unique;not null"` // Unique and non-nullable email
	Password string `gorm:"not null"`        // Non-nullable password (should be hashed)

	// Relationships
	Messages     []Message     `gorm:"foreignKey:SenderID"` // Messages sent by this user
	UserMessages []UserMessage `gorm:"foreignKey:OwnerID"`  // Messages received by this user
	Contacts     []UserContact `gorm:"foreignKey:UserID"`   // Optional, one-to-many contacts
	UserGroups   []UserGroup   `gorm:"foreignKey:UserID"`   // Groups this user belongs to
}

var (
	ErrEmptyPassword = errors.New("password should not be empty")
)

func init() {
	Register(&User{})
}

// HashPassword hashes a plaintext password and returns the hashed string.
func (u *User) HashPassword(plain string) (string, error) {
	if len(plain) == 0 {
		return "", ErrEmptyPassword
	}
	h, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	return string(h), err
}

// CheckPassword verifies if the given plaintext password matches the hashed password.
func (u *User) CheckPassword(plain string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plain))
	return err == nil
}
