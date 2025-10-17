package response

import (
	"strings"

	"github.com/alextilot/golang-htmx-chatapp/constants/header"
	"github.com/alextilot/golang-htmx-chatapp/constants/headerval"
	"github.com/labstack/echo/v4"
)

// wantsJSON checks if the request expects JSON.
// It looks at Accept header and Content-Type header.
func wantsJSON(c echo.Context) bool {
	accept := c.Request().Header.Get(header.Accept)
	ct := c.Request().Header.Get(header.ContentType)
	return strings.Contains(accept, headerval.ContentTypeJSON) || strings.Contains(ct, headerval.ContentTypeJSON)
}
