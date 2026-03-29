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

	Messages []UserMessage // One-to-many relationship: a user can send many messages
	Contacts []UserContact `gorm:"foreignKey:UserID"` // optional, one-to-many
}

var (
	ErrEmptyPassword = errors.New("password should not be empty")
)

func init() {
	Register(&User{})
}

func (u *User) HashPassword(plain string) (string, error) {
	if len(plain) == 0 {
		return "", ErrEmptyPassword
	}
	h, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	return string(h), err
}

func (u *User) CheckPassword(plain string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plain))
	return err == nil
}
