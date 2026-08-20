package ws

import (
	"github.com/alextilot/golang-htmx-chatapp/internal/realtime"
	"github.com/alextilot/golang-htmx-chatapp/internal/service"
)

type Handler struct {
	Services *service.Services
	Hub      *realtime.Hub
	Group    *GroupHandler
}

// NewHandler does not start Hub's event loop — the caller must run
// go handler.Hub.Run(ctx), or every connection will block on registering.
func NewHandler(svc *service.Services) *Handler {
	hub := realtime.NewHub(svc.GroupService)

	return &Handler{
		Services: svc,
		Hub:      hub,
		Group:    NewGroupHandler(hub),
	}
}
