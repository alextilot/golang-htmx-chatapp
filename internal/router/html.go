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
	authenticated := site.Group("")
	authenticated.Use(auth.RequireLogin)

	authenticated.GET("/profile", h.Profile)
	authenticated.GET("/profile/edit", h.ProfileEdit)

	authenticated.GET("/groups", h.Group.List)
	authenticated.POST("/groups", h.Group.Create)
	authenticated.GET("/groups/:groupID", h.Group.Show)
	authenticated.PUT("/groups/:groupID", h.Group.Update)
	authenticated.DELETE("/groups/:groupID", h.Group.Delete)
	authenticated.POST("/groups/:groupID/members", h.Group.AddMember)

	authenticated.GET("/chat", h.Chat.List)
	authenticated.GET("/chat/:id", h.Chat.Room)
	authenticated.POST("/chat/:id/join", h.Chat.Join)
	authenticated.POST("/chat/:id/leave", h.Chat.Leave)
	authenticated.POST("/chat/:id/message", h.Chat.SendMessage)
}
