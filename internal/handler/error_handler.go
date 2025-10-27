package handler

import (
	"github.com/a-h/templ"
	"github.com/alextilot/golang-htmx-chatapp/internal/response"
	"github.com/alextilot/golang-htmx-chatapp/internal/validation"
	"github.com/alextilot/golang-htmx-chatapp/web/pages/status"
	"github.com/labstack/echo/v4"
	"net/http"
)

// Define the custom error handler.
func (h *Handler) HTTPErrorHandler(err error, c echo.Context) {
	// Don't do anything if the response is already committed.
	if c.Response().Committed {
		return
	}

	// Extract the HTTP error code.
	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}

	var fe validation.FieldErrors
	var tmpl templ.Component

	switch code {
	case http.StatusNotFound:
		fe = validation.FieldErrors{validation.FieldRequest: {"Not Found"}}
		tmpl = status.NotFoundPage()
	case http.StatusInternalServerError:
		fe = validation.FieldErrors{validation.FieldServer: {"Internal Server Error"}}
		tmpl = status.ServerErrorPage()
	default:
		fe = validation.FieldErrors{validation.FieldServer: {"An unexpected error occurred"}}
		tmpl = status.ServerErrorPage()
	}

	// IMPORTANT: return the result of Send to avoid "superfluous WriteHeader"
	if err := response.Send(c, response.Response{
		Status: code,
		View:   tmpl,
		Errors: fe,
	}); err != nil {
		c.Logger().Errorf("Failed to send error response: %v", err)
	}
}
