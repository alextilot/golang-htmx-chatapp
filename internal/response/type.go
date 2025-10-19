package response

import "github.com/alextilot/golang-htmx-chatapp/internal/validation"

// Response is what handlers return to be rendered by Respond()
type Response struct {
	Status       int                    // HTTP status code
	Data         any                    // JSON data or other
	Errors       validation.FieldErrors // FieldErrors for JSON or HTML
	HTMLTemplate any                    // templ.Component for HTML
	HTMXRedirect string                 // optional URL for HTMX redirect
}
