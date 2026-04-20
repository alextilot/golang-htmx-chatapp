package router

import (
	"context"
	"github.com/alextilot/golang-htmx-chatapp/internal/handler"
	"github.com/alextilot/golang-htmx-chatapp/internal/routes"
	"github.com/alextilot/golang-htmx-chatapp/web"
	"github.com/alextilot/golang-htmx-chatapp/web/pages"
	"github.com/alextilot/golang-htmx-chatapp/web/pages/status"
	"github.com/labstack/echo/v4"
	"net/http"
)

func RegisterPublicRoutes(e *echo.Echo, h *handler.Handler, ctx context.Context) {
	// Public routes group
	public := e.Group("")

	// Index / landing page
	public.GET(routes.Routes.HomePage.Path, func(c echo.Context) error {
		return web.Render(c, http.StatusOK, pages.HomePage())
	})

	// About page
	public.GET(routes.Routes.AboutPage.Path, func(c echo.Context) error {
		return web.Render(c, http.StatusOK, pages.AboutPage())
	})

	// 401 page
	public.GET(routes.Routes.Unauthorized.Path, func(c echo.Context) error {
		return web.Render(c, http.StatusUnauthorized, status.UnauthorizedPage())
	})

	// 403 page
	public.GET(routes.Routes.Forbidden.Path, func(c echo.Context) error {
		return web.Render(c, http.StatusForbidden, status.ForbiddenPage())
	})

	// 404 page
	public.GET(routes.Routes.NotFound.Path, func(c echo.Context) error {
		return web.Render(c, http.StatusNotFound, status.NotFoundPage())
	})

	// 500 page
	public.GET(routes.Routes.ServerError.Path, func(c echo.Context) error {
		return web.Render(c, http.StatusInternalServerError, status.ServerErrorPage())
	})

	// Sign up
	public.GET(routes.Routes.SignupPage.Path, func(c echo.Context) error {
		return web.Render(c, http.StatusOK, pages.SignupPage())
	})
	public.POST(routes.Routes.SignupPage.Path, h.SignUp)

	// Login
	public.GET(routes.Routes.LoginPage.Path, func(c echo.Context) error {
		return web.Render(c, http.StatusOK, pages.LoginPage())
	})
	public.POST(routes.Routes.LoginPage.Path, h.Login)

	// Logout
	public.POST(routes.Routes.Logout.Path, h.Logout)
}
