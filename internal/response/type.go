package response

import (
	"github.com/a-h/templ"
	"github.com/alextilot/golang-htmx-chatapp/internal/validation"
)

// Response is what handlers return to be rendered by Respond()
type Response struct {
	Status   int                    // HTTP status code
	Redirect string                 // Url to redirect to (HTMX-aware)
	Data     any                    // JSON data or other
	Errors   validation.FieldErrors // Unified errors for JSON or HTML
	View     templ.Component        // templ.Component for HTML
}

// Result describes what happened in a handler.
// It is format-agnostic — the negotiator decides how to render it.
type Result struct {
	Status   int    // HTTP status code
	Payload  any    // JSON payload (success or error body)
	Redirect string // if set, all renderers respect this
}
