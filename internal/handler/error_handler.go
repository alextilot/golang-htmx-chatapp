package handler

import (
	"net/http"
	"time"

	"github.com/a-h/templ"
	"github.com/alextilot/golang-htmx-chatapp/internal/response"
	"github.com/alextilot/golang-htmx-chatapp/internal/validation"
	"github.com/alextilot/golang-htmx-chatapp/web"
	"github.com/alextilot/golang-htmx-chatapp/web/pages/status"
	"github.com/labstack/echo/v5"
)

// Define the custom error handler.
func (h *Handler) HTTPErrorHandler(err error, c *echo.Context) {
	// Skip if response is already committed
	if resp, _ := echo.UnwrapResponse(c.Response()); resp != nil && resp.Committed {
		return
	}

	// Default to internal server error
	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}
	c.Logger().Error("HTTP error", "code", code, "error", err)

	var tmpl templ.Component
	var fe validation.FieldErrors

	switch code {
	case http.StatusNotFound:
		fe = validation.FieldErrors{validation.FieldRequest: {"Page Not Found"}}
		tmpl = status.NotFoundPage()
	case http.StatusForbidden:
		fe = validation.FieldErrors{validation.FieldRequest: {"Forbidden"}}
		tmpl = status.ForbiddenPage()
	case http.StatusUnauthorized:
		fe = validation.FieldErrors{validation.FieldRequest: {"Unauthorized"}}
		tmpl = status.UnauthorizedPage()
	case http.StatusInternalServerError:
		fe = validation.FieldErrors{validation.FieldServer: {"Internal Server Error"}}
		tmpl = status.ServerErrorPage()
	default:
		fe = validation.FieldErrors{validation.FieldServer: {"An unexpected error occurred"}}
		tmpl = status.ServerErrorPage()
	}

	time.Sleep(2 * time.Second)
	web.Render(c, 200, status.ServerErrorPage())
	return

	c.Logger().Debug("Sending error response", "status", code)
	if err := response.Send(c, response.Response{
		Status: code,
		View:   tmpl,
		Errors: fe,
	}); err != nil {
		c.Logger().Error("Failed to send error response", "error", err)

	}
}
