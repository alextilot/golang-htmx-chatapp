package response

import (
	"net/http"

	"github.com/alextilot/golang-htmx-chatapp/internal/constants/header"
	"github.com/alextilot/golang-htmx-chatapp/web"
	"github.com/labstack/echo/v4"
)

// Send negotiates the correct response format based on request headers
// and renders accordingly.
//
// Priority:
//  1. Redirect (HTMX-aware)
//  2. HTMX partial (Fragment)
//  3. HTMX boosted or fresh browser (Page)
//  4. JSON
//  5. NoContent fallback
func Send(c echo.Context, r Result, opts ...Option) error {
	status := r.Status
	if status == 0 {
		status = http.StatusOK
	}

	kind := negotiate(c)

	// --- Redirect ---
	// Handled before format negotiation because it overrides rendering.
	// Each client type gets the redirect in the appropriate way.
	if r.Redirect != "" {
		switch kind {
		case kindHTMX, kindBoosted:
			c.Response().Header().Set(header.HXRedirect, r.Redirect)
			return c.NoContent(http.StatusOK)
		case kindJSON:
			return c.JSON(http.StatusOK, map[string]any{
				"redirect": r.Redirect,
			})
		default:
			return c.Redirect(http.StatusSeeOther, r.Redirect)
		}
	}

	// --- Format negotiation ---
	switch kind {

	case kindHTMX:
		html := findHTML(opts)
		if html == nil {
			break // fall through to JSON
		}
		view := html.fragment
		if view == nil {
			view = html.page // fragment not provided, fall back to full page
		}
		return web.Render(c, status, view)

	case kindBoosted, kindHTML:
		html := findHTML(opts)
		if html == nil {
			break // fall through to JSON
		}
		if html.page == nil {
			break
		}
		return web.Render(c, status, html.page)

	case kindJSON:
		if r.Payload == nil {
			return c.NoContent(status)
		}
		return c.JSON(status, r.Payload)
	}

	// --- Fallback ---
	// No matching renderer found for the requested format.
	// Serialize Payload as JSON if available, otherwise 204.
	if r.Payload != nil {
		return c.JSON(status, r.Payload)
	}
	return c.NoContent(status)
}
