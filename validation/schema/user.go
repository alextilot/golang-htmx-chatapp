package schema

import (
	"strings"

	"github.com/alextilot/golang-htmx-chatapp/validation"
)

type UserLoginInput struct {
	Username string `form:"username" json:"username" validate:"required"`
	Password string `form:"password" json:"password" validate:"required"`
}

func SanitizeUserLogin(input UserLoginInput) UserLoginInput {
	input.Username = strings.TrimSpace(input.Username)

	return input
}

func ValidateUserLogin(input UserLoginInput) validation.FieldErrors {
	return validation.ValidateStruct(input)
}

type UserCreateInput struct {
	Username       string `form:"username" json:"username" validate:"required,min=2,max=20"`
	Email          string `form:"email" json:"email" validate:"required,email"`
	Password       string `form:"password" json:"password" validate:"required,password"`
	RepeatPassword string `form:"repeatPassword" json:"repeatPassword" validate:"required,eqfield=Password"`
}

func SanitizeUserCreate(input UserCreateInput) UserCreateInput {
	input.Username = strings.TrimSpace(input.Username)
	input.Email = strings.TrimSpace(strings.ToLower(input.Email))
	return input
}

func ValidateUserCreate(input UserCreateInput) validation.FieldErrors {
	errors := validation.ValidateStruct(input)

	if input.Password == input.Username {
		errors["Password"] = append(errors["Password"], "Password cannot match username")
	}

	// Check for reserved words
	var reservedUsernames = []string{"admin", "root", "system"}

	for _, word := range reservedUsernames {
		if strings.Contains(strings.ToLower(input.Username), word) {
			errors["Username"] = append(errors["Username"], "Username cannot contain reserved words")
		}
	}

	return errors
}
