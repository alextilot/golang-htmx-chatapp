package response

import (
	"github.com/a-h/templ"
	"github.com/alextilot/golang-htmx-chatapp/internal/validation"
)

// Response is what handlers return to be rendered by Respond()
type Response struct {
	Status   int                    // HTTP status code
	Data     any                    // JSON data or other
	View     templ.Component        // templ.Component for HTML
	Redirect string                 // Url to redirect to (HTMX-aware)
	Errors   validation.FieldErrors // Unified errors for JSON or HTML
}
