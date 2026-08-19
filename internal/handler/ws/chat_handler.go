package ws

import (
	"context"

	"github.com/alextilot/golang-htmx-chatapp/internal/realtime"
	"github.com/labstack/echo/v5"
)

// ChatHandler has no service dependency of its own — a chat room is a
// Group, and membership is enforced by GroupService.SendMessage on every
// inbound message via the hub.
type ChatHandler struct {
	hub *realtime.Hub
}

func NewChatHandler(hub *realtime.Hub) *ChatHandler {
	return &ChatHandler{hub: hub}
}

// ConnectWS upgrades the connection to WebSocket for a chat room.
// GET /ws/chat/:id
func (h *ChatHandler) ConnectWS(c *echo.Context, ctx context.Context) error {
	return connect(c, ctx, h.hub, c.Param("id"))
}
