package schema

import (
	"strings"

	"github.com/alextilot/golang-htmx-chatapp/internal/validation"
)

type ChatMessageInput struct {
	Content string `form:"Content" json:"content" validate:"required,max=5000"`
}

func SanitizeChatMessage(input ChatMessageInput) ChatMessageInput {
	input.Content = strings.TrimSpace(input.Content)
	return input
}

func ValidateChatMessage(input ChatMessageInput) validation.FieldErrors {
	return validation.ValidateStruct(input)
}
