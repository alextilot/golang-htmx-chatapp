package validation

import "github.com/alextilot/golang-htmx-chatapp/internal/validation/core"

// Re-export types
type FieldErrors = core.FieldErrors

// Re-export functions
var (
	Validator            = core.Validator
	ValidateStruct       = core.ValidateStruct
	RegisterErrorMessage = core.RegisterErrorMessage
)

const (
	// FieldServer represents server or system errors (HTTP 500)
	FieldServer = "server"

	// FieldRequest represents input parsing or sanitization errors (HTTP 400)
	FieldRequest = "request"

	// FieldValidation represents input content or business validation errors (HTTP 400)
	FieldValidation = "validation"

	// FieldAuth represents authentication or authorization errors (HTTP 401 / 403)
	FieldAuth = "auth"

	// FieldGlobal represents global or UI messages, e.g., form-level notices
	FieldGlobal = "global"
)

func NewFieldErrors() FieldErrors {
	return core.NewFieldErrors()
}
