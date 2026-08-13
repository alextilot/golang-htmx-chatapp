package router

import (
	"context"
	"strings"
	"time"

	"github.com/alextilot/golang-htmx-chatapp/internal/auth"
	"github.com/alextilot/golang-htmx-chatapp/internal/config"
	"github.com/alextilot/golang-htmx-chatapp/internal/handler"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

// NewRouter creates and configures the application's Echo server.
func NewRouter(ctx context.Context, h *handler.Handlers) *echo.Echo {
	e := echo.New()

	e.Pre(middleware.RemoveTrailingSlash())

	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())
	e.Use(middleware.Gzip())
	e.Use(middleware.Secure())

	e.Use(middleware.ContextTimeoutWithConfig(
		middleware.ContextTimeoutConfig{
			Timeout: 10 * time.Second,
		},
	))

	e.Use(auth.AuthMiddleware)
	e.Use(cacheControlMiddleware)

	e.Static("/static", "web/static")

	registerHTMLRoutes(e, h.HTML)
	registerAPIRoutes(e, h.API)
	registerWebSocketRoutes(e, ctx, h.WS)

	return e
}

func cacheControlMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		path := c.Request().URL.Path

		if config.Cfg.IsProduction() && strings.HasPrefix(path, "/static/") {
			c.Response().Header().Set(
				"Cache-Control",
				"public, max-age=2592000, immutable",
			)

			return next(c)
		}

		c.Response().Header().Set(
			"Cache-Control",
			"no-store, no-cache, max-age=0, must-revalidate",
		)
		c.Response().Header().Set("Pragma", "no-cache")
		c.Response().Header().Set("Expires", "0")

		return next(c)
	}
}
