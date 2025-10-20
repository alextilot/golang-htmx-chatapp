package router

import (
	"github.com/alextilot/golang-htmx-chatapp/internal/auth"
	"github.com/alextilot/golang-htmx-chatapp/internal/handler"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"time"
)

func NewRouter(h *handler.Handler) *echo.Echo {
	e := echo.New()

	e.Pre(middleware.RemoveTrailingSlash())

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.Gzip())
	e.Use(middleware.Secure())
	// e.Use(middleware.CORS())
	// e.Use(middleware.CSRF())
	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Timeout: 10 * time.Second,
	}))

	e.Use(auth.TokenRefresherMiddleware)
	e.Use(auth.UserContextMiddleware)

	// Serve static files
	e.Static("/static", "web/static")

	RegisterPublicRoutes(e, h)
	RegisterPrivateRoutes(e, h)

	return e
}
