package web

import (
	"strings"

	"github.com/labstack/echo/v4"
)

// wantsJSON checks if the request expects JSON.
// It looks at Accept header and Content-Type header.
func wantsJSON(c echo.Context) bool {
	accept := c.Request().Header.Get("Accept")
	ct := c.Request().Header.Get("Content-Type")
	return strings.Contains(accept, "application/json") || strings.Contains(ct, "application/json")
}
