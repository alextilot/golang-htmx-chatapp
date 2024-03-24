package model

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Username string
	Password string
}

var (
	ErrEmptyPassword = errors.New("password should not be empty")
)

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
