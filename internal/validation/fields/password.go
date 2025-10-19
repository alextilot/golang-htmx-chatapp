package fields

import (
	"regexp"

	"github.com/go-playground/validator/v10"

	"github.com/alextilot/golang-htmx-chatapp/internal/validation"
)

func init() {
	validation.Validator().RegisterValidation("password", func(fl validator.FieldLevel) bool {
		password := fl.Field().String()
		var (
			hasMinLen  = len(password) >= 8
			hasUpper   = regexp.MustCompile(`[A-Z]`).MatchString(password)
			hasLower   = regexp.MustCompile(`[a-z]`).MatchString(password)
			hasNumber  = regexp.MustCompile(`[0-9]`).MatchString(password)
			hasSpecial = regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`).MatchString(password)
		)
		return hasMinLen && hasUpper && hasLower && hasNumber && hasSpecial
	})

	// Register the default error message for password
	validation.RegisterErrorMessage("password", "Password must be at least 8 characters and include an uppercase letter, a lowercase letter, a number, and a special character.")
}
