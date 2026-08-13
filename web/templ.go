package web

import (
	"fmt"
	"net/http"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v5"
)

func Render(ctx *echo.Context, status int, t templ.Component) error {
	buf := templ.GetBuffer()
	defer templ.ReleaseBuffer(buf)

	if err := t.Render(ctx.Request().Context(), buf); err != nil {
		if ctx.Request().Context().Err() != nil {
			ctx.Logger().Warn("Render aborted: context canceled")
			return nil
		}
		ctx.Logger().Error("Failed to render template", "type", fmt.Sprintf("%T", t), "error", err)
		return ctx.String(http.StatusInternalServerError, "Failed to render page")
	}

	return ctx.HTML(status, buf.String())
}
