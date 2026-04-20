package response

import (
	"strings"

	"github.com/alextilot/golang-htmx-chatapp/internal/constants/header"
	"github.com/labstack/echo/v4"
)

// wantsJSON checks if the request expects JSON.
// It looks at Accept header and Content-Type header.
func wantsJSON(c echo.Context) bool {
	accept := c.Request().Header.Get(echo.HeaderAccept)
	ct := c.Request().Header.Get(echo.HeaderContentType)
	return strings.Contains(accept, echo.HeaderContentType) || strings.Contains(ct, echo.MIMEApplicationJSON)
}

func isHTMX(c echo.Context) bool {
	return c.Request().Header.Get(header.HXRequest) == "true"
}
