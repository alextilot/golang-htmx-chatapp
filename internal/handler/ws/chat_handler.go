package ws

import (
	"context"

	"github.com/alextilot/golang-htmx-chatapp/internal/service"
	"github.com/labstack/echo/v5"
)

// ChatHandler handles WebSocket connections for chat rooms.
//
// A chat room is a Group (see service.GroupService) — it currently only
// upgrades and registers the connection with the hub; the service is
// retained for future authorization checks (e.g. verifying membership
// before allowing a connection).
type ChatHandler struct {
	hub *Hub
	svc *service.GroupService
}

func NewChatHandler(hub *Hub, svc *service.GroupService) *ChatHandler {
	return &ChatHandler{hub: hub, svc: svc}
}

// ConnectWS upgrades the connection to WebSocket for a chat room.
// GET /ws/chat/:id
func (h *ChatHandler) ConnectWS(c *echo.Context, ctx context.Context) error {
	return h.hub.Handler(c, ctx, c.Param("id"))
}
