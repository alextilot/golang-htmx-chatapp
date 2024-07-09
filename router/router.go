package router

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func New() *echo.Echo {
	e := echo.New()
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Recover())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))
	e.Use(globalResponseHeader)
	e.Validator = NewValidator()
	return e
}

// TODO: with the forced redirects from /chatroom to /login for unauth users
// the browser is/was caching the page wrong and not allowing /chatroom page to load.
// Not 100% sure if it is fixed with Cache-control set but I havn't seen the error.
func globalResponseHeader(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set("Cache-Control", "max-age=0, no-cache, no-store")
		return next(c)
	}
}
