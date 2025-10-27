package response

import (
	"github.com/alextilot/golang-htmx-chatapp/internal/constants/header"
	"github.com/alextilot/golang-htmx-chatapp/web"
	"github.com/labstack/echo/v4"
	"net/http"
)

func Send(c echo.Context, r Response) error {
	status := r.Status
	if status == 0 {
		status = http.StatusOK
	}

	// Handle HTMX redirect
	if r.Redirect != "" && isHTMX(c) {
		c.Response().Header().Set(header.HXRedirect, r.Redirect)
		return c.NoContent(http.StatusOK)
	}

	// Handle redirect (non-HTMX)
	if r.Redirect != "" {
		return c.Redirect(http.StatusSeeOther, r.Redirect)
	}

	// Handle JSON
	if wantsJSON(c) {
		jsonData := r.Data

		if r.Errors != nil {
			jsonData = map[string]any{"errors": r.Errors.FlattenByField()}
		}

		if jsonData == nil {
			jsonData = map[string]any{}
		}

		return c.JSON(status, jsonData)
	}

	// Handle HTML template
	if r.View != nil {
		return web.Render(c, status, r.View)
	}

	// Fallback
	return c.NoContent(status)
}
