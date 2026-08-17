package router

import (
	"net/http"

	"github.com/alextilot/golang-htmx-chatapp/internal/auth"
	"github.com/alextilot/golang-htmx-chatapp/internal/handler/html"
	"github.com/alextilot/golang-htmx-chatapp/web"
	"github.com/alextilot/golang-htmx-chatapp/web/pages"
	"github.com/alextilot/golang-htmx-chatapp/web/pages/status"

	"github.com/labstack/echo/v5"
)

func registerHTMLRoutes(e *echo.Echo, h *html.Handler) {
	site := e.Group("")

	// Public website.
	site.GET("/", func(c *echo.Context) error {
		return web.Render(c, http.StatusOK, pages.HomePage())
	})

	site.GET("/about", func(c *echo.Context) error {
		return web.Render(c, http.StatusOK, pages.AboutPage())
	})

	site.GET("/signup", func(c *echo.Context) error {
		return web.Render(c, http.StatusOK, pages.SignupPage())
	})

	site.POST("/signup", h.SignUp)

	site.GET("/login", func(c *echo.Context) error {
		return web.Render(c, http.StatusOK, pages.LoginPage())
	})

	site.POST("/login", h.Login)

	site.POST("/logout", h.Logout)

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
	site.GET("/profile", h.Profile, auth.RequireLogin)
	site.GET("/profile/edit", h.ProfileEdit, auth.RequireLogin)

	site.GET("/groups", h.Group.List, auth.RequireLogin)
	site.POST("/groups", h.Group.Create, auth.RequireLogin)
	site.GET("/groups/:groupID", h.Group.Show, auth.RequireLogin)
	site.PUT("/groups/:groupID", h.Group.Update, auth.RequireLogin)
	site.DELETE("/groups/:groupID", h.Group.Delete, auth.RequireLogin)
	site.POST("/groups/:groupID/members", h.Group.AddMember, auth.RequireLogin)

	site.GET("/chat", h.Chat.List, auth.RequireLogin)
	site.GET("/chat/:id", h.Chat.Room, auth.RequireLogin)
	site.POST("/chat/:id/join", h.Chat.Join, auth.RequireLogin)
	site.POST("/chat/:id/leave", h.Chat.Leave, auth.RequireLogin)
	site.POST("/chat/:id/message", h.Chat.SendMessage, auth.RequireLogin)
}
