package html

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/alextilot/golang-htmx-chatapp/web"
	"github.com/alextilot/golang-htmx-chatapp/web/pages/status"
	"github.com/labstack/echo/v5"
)

// HTTPErrorHandler is Echo's global fallback error handler for the HTML
// website. It renders an HTML status page for errors that reach Echo
// without having already been mapped to a response by a handler.
//
// API and WebSocket errors are mapped by their own handlers and never reach
// this function; it is only registered for the website's error path.
func (h *Handler) HTTPErrorHandler(err error, c *echo.Context) {
	// Skip if the response is already committed.
	if resp, _ := echo.UnwrapResponse(c.Response()); resp != nil && resp.Committed {
		return
	}

	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}
	c.Logger().Error("HTTP error", "code", code, "error", err)

	var page templ.Component

	switch code {
	case http.StatusNotFound:
		page = status.NotFoundPage()
	case http.StatusForbidden:
		page = status.ForbiddenPage()
	case http.StatusUnauthorized:
		page = status.UnauthorizedPage()
	default:
		page = status.ServerErrorPage()
	}

	c.Logger().Debug("Sending error response", "status", code)
	if err := web.Render(c, code, page); err != nil {
		c.Logger().Error("Failed to send error response", "error", err)
	}
}
