package router

import (
	"context"
	"strings"
	"time"

	"github.com/alextilot/golang-htmx-chatapp/internal/auth"
	"github.com/alextilot/golang-htmx-chatapp/internal/handler"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

// Deps holds the dependencies needed to construct and configure the
// application's Echo server.
type Deps struct {
	Ctx      context.Context
	Handlers *handler.Handlers
	AuthSvc  *auth.Service
	// IsStaticCacheEnabled controls whether /static assets get long-lived
	// cache headers. It's derived once, in cmd/main, from the environment
	// rather than read from a global config here.
	IsStaticCacheEnabled bool
}

// NewRouter creates and configures the application's Echo server.
func NewRouter(deps Deps) *echo.Echo {
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

	e.Use(deps.AuthSvc.AuthMiddleware)
	e.Use(cacheControlMiddleware(deps.IsStaticCacheEnabled))

	e.HTTPErrorHandler = deps.Handlers.HTML.HTTPErrorHandler

	e.Static("/static", "web/static")

	registerHTMLRoutes(e, deps.Handlers.HTML)
	registerAPIRoutes(e, deps.Handlers.API)
	registerWebSocketRoutes(e, deps.Ctx, deps.Handlers.WS)

	return e
}

func cacheControlMiddleware(isStaticCacheEnabled bool) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			path := c.Request().URL.Path

			if isStaticCacheEnabled && strings.HasPrefix(path, "/static/") {
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
}
