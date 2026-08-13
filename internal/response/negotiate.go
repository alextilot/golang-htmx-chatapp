package response

import (
	"strings"

	"github.com/alextilot/golang-htmx-chatapp/internal/constants/header"
	"github.com/labstack/echo/v5"
)

type requestKind int

const (
	kindJSON    requestKind = iota
	kindHTMX                // HX-Request: true, HX-Boosted: false
	kindBoosted             // HX-Request: true, HX-Boosted: true
	kindHTML                // fresh browser load
)

// negotiate inspects request headers and returns the requestKind.
// Priority: HTMX > JSON > HTML
func negotiate(c *echo.Context) requestKind {
	if isHTMX(c) {
		if isBoosted(c) {
			return kindBoosted
		}
		return kindHTMX
	}
	if wantsJSON(c) {
		return kindJSON
	}
	return kindHTML
}

func isHTMX(c *echo.Context) bool {
	return c.Request().Header.Get(header.HXRequest) == "true"
}

func isBoosted(c *echo.Context) bool {
	return c.Request().Header.Get(header.HXBoosted) == "true"
}

func wantsJSON(c *echo.Context) bool {
	accept := c.Request().Header.Get(echo.HeaderAccept)
	ct := c.Request().Header.Get(echo.HeaderContentType)
	return strings.Contains(accept, echo.MIMEApplicationJSON) ||
		strings.Contains(ct, echo.MIMEApplicationJSON)
}

func wantsHTML(c *echo.Context) bool {
	accept := c.Request().Header.Get(echo.HeaderAccept)
	return strings.Contains(accept, echo.MIMETextHTML)
}

// findHTML returns the htmlOption from opts, or nil if not registered.
func findHTML(opts []Option) *htmlOption {
	for _, o := range opts {
		if h, ok := o.(htmlOption); ok {
			return &h
		}
	}
	return nil
}
