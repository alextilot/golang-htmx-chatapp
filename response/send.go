package response

import (
	"github.com/a-h/templ"
	"github.com/alextilot/golang-htmx-chatapp/constants/header"
	"github.com/alextilot/golang-htmx-chatapp/web"
	"github.com/labstack/echo/v4"
	"net/http"
)

func Send(c echo.Context, resp Response) error {
	status := resp.Status
	if status == 0 {
		status = http.StatusOK
	}

	// HTMX redirect
	if resp.HTMXRedirect != "" && c.Request().Header.Get(header.HXRedirect) == "true" {
		c.Response().Header().Set(header.HXRedirect, resp.HTMXRedirect)
		return c.NoContent(http.StatusOK)
	}

	// JSON
	if wantsJSON(c) {
		jsonData := resp.Data

		if resp.Errors != nil {
			jsonData = map[string]any{"errors": resp.Errors.FlattenByField()}
		}

		if jsonData == nil {
			jsonData = map[string]any{}
		}

		return c.JSON(status, jsonData)
	}

	// HTML
	if tmpl, ok := resp.HTMLTemplate.(templ.Component); ok {
		return web.Render(c, status, tmpl)
	}

	return c.NoContent(status)
}
