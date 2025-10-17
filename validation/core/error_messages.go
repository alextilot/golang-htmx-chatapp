package core

import (
	"strings"
)

// Error message templates by tag name.
// Use {{field}} and {{param}} placeholders for flexibility.
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
