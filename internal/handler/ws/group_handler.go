package ws

import (
	"context"

	"github.com/alextilot/golang-htmx-chatapp/internal/service"
	"github.com/labstack/echo/v5"
)

// GroupHandler handles WebSocket connections for group chat.
//
// It currently only upgrades and registers the connection with the hub;
// the group service is retained for future authorization checks (e.g.
// verifying membership before allowing a connection).
type GroupHandler struct {
	hub *Hub
	svc *service.GroupService
}

func NewGroupHandler(hub *Hub, svc *service.GroupService) *GroupHandler {
	return &GroupHandler{hub: hub, svc: svc}
}

// ConnectWS upgrades the connection to WebSocket for a group's chat.
// GET /ws/groups/:groupID
func (h *GroupHandler) ConnectWS(c *echo.Context, ctx context.Context) error {
	return h.hub.Handler(c, ctx, c.Param("groupID"))
}
