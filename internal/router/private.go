package router

import (
	"context"

	"github.com/alextilot/golang-htmx-chatapp/internal/auth"
	"github.com/alextilot/golang-htmx-chatapp/internal/handler"
	"github.com/labstack/echo/v4"
)

func RegisterPrivateRoutes(e *echo.Echo, h *handler.Handler, ctx context.Context) {
	private := e.Group("/app")
	private.Use(auth.RequireLogin)

	// ---- Groups ----
	private.GET("/groups", h.Group.List)
	private.POST("/groups", h.Group.Create)
	private.GET("/groups/:groupID", h.Group.Show)
	private.PUT("/groups/:groupID", h.Group.Update)
	private.DELETE("/groups/:groupID", h.Group.Delete)
	private.POST("/groups/:groupID/members", h.Group.AddMember)

	// ---- WebSocket ----
	// Sits outside /app so the WS upgrade is not affected by session middleware.
	e.GET("/ws/groups/:groupID", func(c echo.Context) error {
		return h.Group.ConnectWS(c, ctx)
	})
}
