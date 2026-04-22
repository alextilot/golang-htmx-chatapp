package router

import (
	"context"
	"github.com/alextilot/golang-htmx-chatapp/internal/auth"
	"github.com/alextilot/golang-htmx-chatapp/internal/config"
	"github.com/alextilot/golang-htmx-chatapp/internal/handler"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"strings"
	"time"
)

// CacheControlMiddleware sets cache headers intelligently:
// - Static assets: cached for 30 days
// - Dynamic pages/API: no cache
func cacheControlMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		path := c.Request().URL.Path

		if config.Cfg.IsProduction() && strings.HasPrefix(path, "/static/") {
			// 30 days cache for static assets in production
			c.Response().Header().Set("Cache-Control", "public, max-age=2592000, immutable")
			return next(c)
		}

		// Everything else (dynamic pages, or non-prod static) → no cache
		c.Response().Header().Set("Cache-Control", "no-store, no-cache, max-age=0, must-revalidate")
		c.Response().Header().Set("Pragma", "no-cache")
		c.Response().Header().Set("Expires", "0")

		return next(c)
	}
}

func NewRouter(h *handler.Handler, ctx context.Context) *echo.Echo {
	e := echo.New()

	e.Debug = config.Cfg.Debug
	e.HTTPErrorHandler = h.HTTPErrorHandler

	// ---- Pre middleware ----
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))

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
	e.Use(auth.AuthMiddleware)

	e.Use(cacheControlMiddleware)

	// ---- Static assets ----
	e.Static("/static", "web/static")

	// ---- Routes ----
	RegisterPublicRoutes(e, h, ctx)
	RegisterPrivateRoutes(e, h, ctx)

	return e
}
