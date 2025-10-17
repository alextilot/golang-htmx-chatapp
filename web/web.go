package web

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func Respond(c echo.Context, resp Response) error {
	status := resp.Status
	if status == 0 {
		status = http.StatusOK
	}

	// 1️⃣ HTMX redirect
	if resp.HTMXRedirect != "" && c.Request().Header.Get("HX-Request") == "true" {
		c.Response().Header().Set("HX-Redirect", resp.HTMXRedirect)
		return c.NoContent(http.StatusOK)
	}

	// 2️⃣ JSON
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

	// 3️⃣ HTML
	if tmpl, ok := resp.HTMLTemplate.(templ.Component); ok {
		return Render(c, status, tmpl)
	}

	// fallback
	return c.NoContent(status)
}
