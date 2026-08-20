package router

import (
	"net/http"

	"github.com/alextilot/golang-htmx-chatapp/internal/auth"
	"github.com/alextilot/golang-htmx-chatapp/internal/handler/html"
	"github.com/alextilot/golang-htmx-chatapp/internal/routes"
	"github.com/alextilot/golang-htmx-chatapp/web"
	"github.com/alextilot/golang-htmx-chatapp/web/pages"
	"github.com/alextilot/golang-htmx-chatapp/web/pages/status"

	"github.com/labstack/echo/v5"
)

func registerHTMLRoutes(e *echo.Echo, h *html.Handler) {
	site := e.Group("")

	// Public website.
	site.GET(routes.Routes.HomePage.Path, func(c *echo.Context) error {
		return web.Render(c, http.StatusOK, pages.HomePage())
	})

	site.GET(routes.Routes.AboutPage.Path, func(c *echo.Context) error {
		return web.Render(c, http.StatusOK, pages.AboutPage())
	})

	site.GET(routes.Routes.SignupPage.Path, func(c *echo.Context) error {
		return web.Render(c, http.StatusOK, pages.SignupPage())
	})

	site.POST(routes.Routes.SignupPage.Path, h.SignUp)

	site.GET(routes.Routes.LoginPage.Path, func(c *echo.Context) error {
		return web.Render(c, http.StatusOK, pages.LoginPage())
	})

	site.POST(routes.Routes.LoginPage.Path, h.Login)

	site.POST(routes.Routes.Logout.Path, h.Logout)

	// Error pages.
	site.GET("/401", func(c *echo.Context) error {
		return web.Render(
			c,
			http.StatusUnauthorized,
			status.UnauthorizedPage(),
		)
	})

	site.GET("/403", func(c *echo.Context) error {
		return web.Render(
			c,
			http.StatusForbidden,
			status.ForbiddenPage(),
		)
	})

	site.GET("/404", func(c *echo.Context) error {
		return web.Render(
			c,
			http.StatusNotFound,
			status.NotFoundPage(),
		)
	})

	site.GET("/500", func(c *echo.Context) error {
		return web.Render(
			c,
			http.StatusInternalServerError,
			status.ServerErrorPage(),
		)
	})

	// Authenticated website.
	//
	// RequireLogin is applied per-route rather than via a group + Use(), because
	// Echo v5's Group.Use() auto-registers a "/*" catch-all 404 route carrying the
	// group's middleware so it still runs on non-matching paths. This group shares
	// the site group's empty prefix, so that catch-all would intercept every
	// unmatched path on the whole site and turn 404s into 401s.
	site.GET(routes.Routes.Profile.Path, h.Profile, auth.RequireLogin)
	site.GET(routes.Routes.Profile.Path+"/edit", h.ProfileEdit, auth.RequireLogin)

	site.GET(routes.Routes.Groups.Path, h.Group.List, auth.RequireLogin)
	site.POST(routes.Routes.Groups.Path, h.Group.Create, auth.RequireLogin)
	site.GET(routes.Routes.Groups.Path+"/:groupID", h.Group.Show, auth.RequireLogin)
	site.GET(routes.Routes.Groups.Path+"/:groupID/messages", h.Group.Messages, auth.RequireLogin)
	site.PUT(routes.Routes.Groups.Path+"/:groupID", h.Group.Update, auth.RequireLogin)
	site.DELETE(routes.Routes.Groups.Path+"/:groupID", h.Group.Delete, auth.RequireLogin)
	site.POST(routes.Routes.Groups.Path+"/:groupID/members", h.Group.AddMember, auth.RequireLogin)
}
