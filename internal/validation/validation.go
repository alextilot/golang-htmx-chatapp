// Package validation holds field-format primitives — facts about what a
// field's format must look like, independent of any transport or service.
// It is consumed by exactly one caller, internal/apperr.ValidateStruct; see
// docs/error-handling.md for the ownership rules this package follows.
package validation

import (
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// Validator returns the shared validator instance. Used by
// internal/apperr.ValidateStruct; nothing else should need it directly.
func Validator() *validator.Validate {
	return validate
}

// Error message templates by tag name, used to turn a validator.FieldError
// into a human-readable string. Use {{field}} and {{param}} placeholders.
var errorMessages = map[string]string{
	"required": "The {{field}} field is required.",
	"min":      "The {{field}} must be at least {{param}} characters.",
	"max":      "The {{field}} must be at most {{param}} characters.",
	"email":    "Please enter a valid email address.",
	"eqfield":  "The {{field}} must match {{param}}.",
}

// RegisterErrorMessage allows overriding or adding new templates dynamically.
func RegisterErrorMessage(tag, message string) {
	errorMessages[tag] = message
}

// GetErrorMessage formats a message template with field and param.
func GetErrorMessage(tag, field, param string) string {
	tmpl, ok := errorMessages[tag]
	if !ok {
		tmpl = "{{field}} is invalid."
	}
	msg := strings.ReplaceAll(tmpl, "{{field}}", field)
	msg = strings.ReplaceAll(msg, "{{param}}", param)
	return msg
}

func init() {
	validate.RegisterValidation("password", func(fl validator.FieldLevel) bool {
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

	RegisterErrorMessage("password", "Password must be at least 8 characters and include an uppercase letter, a lowercase letter, a number, and a special character.")
}
