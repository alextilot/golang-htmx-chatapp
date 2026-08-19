package ws

import (
	"context"

	"github.com/alextilot/golang-htmx-chatapp/internal/realtime"
	"github.com/labstack/echo/v5"
)

// GroupHandler has no service dependency of its own — membership is
// enforced by GroupService.SendMessage on every inbound message via the hub.
type GroupHandler struct {
	hub *realtime.Hub
}

func NewGroupHandler(hub *realtime.Hub) *GroupHandler {
	return &GroupHandler{hub: hub}
}

// ConnectWS upgrades the connection to WebSocket for a group's chat.
// GET /ws/groups/:groupID
func (h *GroupHandler) ConnectWS(c *echo.Context, ctx context.Context) error {
	return connect(c, ctx, h.hub, c.Param("groupID"))
}
