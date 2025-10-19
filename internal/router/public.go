package router

import (
	"github.com/alextilot/golang-htmx-chatapp/internal/handler"
	"github.com/alextilot/golang-htmx-chatapp/web"
	"github.com/alextilot/golang-htmx-chatapp/web/pages"
	"github.com/labstack/echo/v4"
	"net/http"
)

func RegisterPublicRoutes(e *echo.Echo, h *handler.Handler) {
	// Public routes
	public := e.Group("/")

	e.GET("", func(etx echo.Context) error {
		return web.Render(etx, http.StatusOK, pages.HomePage())
	})

	public.GET("welcome", func(c echo.Context) error {
		return c.String(200, "Welcome!")
	})

	// public.GET("/signup", h.SignUpForm)
	public.POST("signup", h.SignUp)
	// public.GET("/login", h.LoginForm)
	public.POST("login", h.Login)
}
