package router

import (
	"github.com/alextilot/golang-htmx-chatapp/web"
	"github.com/alextilot/golang-htmx-chatapp/web/static"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"net/http"
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
	e.HTTPErrorHandler = customHTTPErrorHandler
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

func customHTTPErrorHandler(err error, etx echo.Context) {
	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}
	etx.Logger().Error(err)

	page := static.GetHttpErrorPage(code)
	if page == nil {
		etx.Logger().Error("HTTP error page not found")
		etx.JSON(code, err)
		return
	}
	if err := web.Render(etx, http.StatusOK, page); err != nil {
		etx.Logger().Error(err)
	}
}
