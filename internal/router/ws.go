package router

import (
	"context"

	"github.com/alextilot/golang-htmx-chatapp/internal/auth"
	"github.com/alextilot/golang-htmx-chatapp/internal/handler/ws"

	"github.com/labstack/echo/v5"
)

func registerWebSocketRoutes(
	e *echo.Echo,
	ctx context.Context,
	h *ws.Handler,
) {
	authenticated := e.Group("/ws")
	authenticated.Use(auth.RequireLogin)

	authenticated.GET("/groups/:groupID", func(c *echo.Context) error {
		return h.Group.ConnectWS(c, ctx)
	})
}
