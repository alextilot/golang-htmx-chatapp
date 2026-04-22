package web

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func Render(ctx echo.Context, status int, t templ.Component) error {
	buf := templ.GetBuffer()
	defer templ.ReleaseBuffer(buf)

	if err := t.Render(ctx.Request().Context(), buf); err != nil {
		if ctx.Request().Context().Err() != nil {
			ctx.Logger().Warn("Render aborted: context canceled")
			return nil
		}
		ctx.Logger().Errorf("Failed to render template %T: %v", t, err)
		return ctx.String(http.StatusInternalServerError, "Failed to render page")
	}

	return ctx.HTML(status, buf.String())
}
