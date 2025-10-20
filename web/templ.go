package web

import (
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"net/http"
)

func Render(ctx echo.Context, status int, t templ.Component) error {
	// Get the request context should have UserContext from middleware
	reqCtx := ctx.Request().Context()

	// Set HTTP status first
	ctx.Response().Writer.WriteHeader(status)

	// Render the templ component using the request context
	err := t.Render(reqCtx, ctx.Response().Writer)
	if err != nil {
		return ctx.String(http.StatusInternalServerError, "failed to render response template")
	}

	return nil
}
