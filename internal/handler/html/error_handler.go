package html

import (
	"context"
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
func (h *Handler) HTTPErrorHandler(c *echo.Context, err error) {
	// Skip if the response is already committed.
	if resp, _ := echo.UnwrapResponse(c.Response()); resp != nil && resp.Committed {
		return
	}

	code := echo.StatusCode(err)
	if code == 0 {
		code = http.StatusInternalServerError
	}
	c.Logger().Error("HTTP error", "code", code, "error", err)

	// By the time Echo invokes the global error handler, the middleware chain
	// has fully unwound — including ContextTimeoutWithConfig's deferred
	// cancel(), which has already canceled the request context. Render with a
	// fresh context so the error page isn't silently dropped as "aborted".
	c.SetRequest(c.Request().WithContext(context.Background()))

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
