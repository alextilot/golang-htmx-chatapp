package html

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/alextilot/golang-htmx-chatapp/internal/apperr"
	"github.com/alextilot/golang-htmx-chatapp/internal/constants/hx"
	"github.com/alextilot/golang-htmx-chatapp/web"
	"github.com/labstack/echo/v5"
)

// Response is what HTML handlers build to describe the outcome of a
// request. sendResponse renders it as the appropriate HTML representation
// (full page, HTMX fragment, or a plain redirect/no-content response).
//
// This type, and the rendering behind it, exists so HTML/HTMX-specific
// concerns (page vs. fragment, HX-Redirect, boosted navigation) stay in the
// HTML handler package instead of a shared abstraction that tried to make
// every transport look the same.
type Response struct {
	Status   int                    // HTTP status code
	Redirect string                 // URL to redirect to (HTMX-aware)
	Errors   apperr.FieldErrors // field errors surfaced to the view
	Page     templ.Component        // full page — fresh browser load or HX-Boosted
	Fragment templ.Component        // partial — HTMX swap; falls back to Page if nil
}

type requestKind int

const (
	kindHTMX    requestKind = iota // HX-Request: true, HX-Boosted: false
	kindBoosted                    // HX-Request: true, HX-Boosted: true
	kindHTML                       // fresh browser load
)

// negotiate inspects request headers to decide how the response should be
// rendered for this HTML/HTMX request.
func negotiate(c *echo.Context) requestKind {
	if isHTMX(c) {
		if isBoosted(c) {
			return kindBoosted
		}
		return kindHTMX
	}
	return kindHTML
}

func isHTMX(c *echo.Context) bool {
	return c.Request().Header.Get(hx.HXRequest) == "true"
}

func isBoosted(c *echo.Context) bool {
	return c.Request().Header.Get(hx.HXBoosted) == "true"
}

// sendResponse renders r as an HTML response.
//
// Priority:
//  1. Redirect (HTMX-aware — HX-Redirect header for HTMX/boosted requests,
//     a normal 303 redirect otherwise)
//  2. HTMX partial (Fragment, falling back to Page)
//  3. Boosted or fresh browser navigation (Page)
//  4. No content, if nothing else applies
func sendResponse(c *echo.Context, r Response) error {
	status := r.Status
	if status == 0 {
		status = http.StatusOK
	}

	kind := negotiate(c)

	if r.Redirect != "" {
		switch kind {
		case kindHTMX, kindBoosted:
			c.Response().Header().Set(hx.HXRedirect, r.Redirect)
			return c.NoContent(http.StatusOK)
		default:
			return c.Redirect(http.StatusSeeOther, r.Redirect)
		}
	}

	switch kind {
	case kindHTMX:
		view := r.Fragment
		if view == nil {
			view = r.Page
		}
		if view != nil {
			return web.Render(c, status, view)
		}
	default: // kindBoosted, kindHTML
		if r.Page != nil {
			return web.Render(c, status, r.Page)
		}
	}

	return c.NoContent(status)
}
