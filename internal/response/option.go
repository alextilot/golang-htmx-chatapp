package response

import (
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

// Option is a renderer registered with Send.
// Each Option handles one content type.
type Option interface {
	contentType() string
}

// htmlOption carries the two HTML representations of a result.
type htmlOption struct {
	page     templ.Component // full page — fresh browser load or HX-Boosted
	fragment templ.Component // partial — HTMX swap
}

func (h htmlOption) contentType() string { return echo.MIMETextHTML }

// ForHTML registers HTML renderers for a result.
// page    — rendered for fresh browser requests and HX-Boosted requests.
// fragment — rendered for HTMX partial swap requests.
// If fragment is nil, page is used as the fallback.
func ForHTML(page, fragment templ.Component) Option {
	return htmlOption{page: page, fragment: fragment}
}
