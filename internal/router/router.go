package router

import (
	"github.com/alextilot/golang-htmx-chatapp/internal/auth"
	"github.com/alextilot/golang-htmx-chatapp/internal/config"
	"github.com/alextilot/golang-htmx-chatapp/internal/handler"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"time"
)

func NewRouter(h *handler.Handler) *echo.Echo {
	e := echo.New()

	e.Debug = config.Cfg.Debug

	// ---- Pre middleware ----
	e.Pre(middleware.RemoveTrailingSlash())

	// ---- Global middleware ----
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.Gzip())
	e.Use(middleware.Secure())

	// Add CORS/CSRF when needed
	// e.Use(middleware.CORS())
	// e.Use(middleware.CSRF())

	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Timeout: 10 * time.Second,
	}))

	// ---- Custom middleware ----
	e.Use(auth.UserContextMiddleware)
	e.Use(auth.TokenRefresherMiddleware)

	// ---- Static assets ----
	e.Static("/static", "web/static")

	e.HTTPErrorHandler = h.HTTPErrorHandler

	// ---- Routes ----
	RegisterPublicRoutes(e, h)
	RegisterPrivateRoutes(e, h)

	return e
}
