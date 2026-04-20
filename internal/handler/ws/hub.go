package ws

import (
	"context"
	"net/http"

	"github.com/alextilot/golang-htmx-chatapp/internal/usercontext"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type clientEventKind int

const (
	eventAdd    clientEventKind = iota
	eventRemove clientEventKind = iota
)

type clientEvent struct {
	kind   clientEventKind
	client *Client
}

// Hub maintains the set of active clients and routes messages between them.
// All client-list mutations are serialized through the events channel to
// avoid the need for a mutex.
type Hub struct {
	clients map[string]*Client
	events  chan clientEvent
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// TODO: restrict to your actual origin(s) before going to production.
		return true
	},
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[string]*Client),
		events:  make(chan clientEvent),
	}
}

// Run processes hub events until the context is cancelled.
// Call this in its own goroutine: go hub.Run(ctx)
func (h *Hub) Run(ctx context.Context) {
	for {
		select {
		case ev := <-h.events:
			switch ev.kind {
			case eventAdd:
				h.clients[ev.client.ID] = ev.client
			case eventRemove:
				if _, ok := h.clients[ev.client.ID]; ok {
					delete(h.clients, ev.client.ID)
					close(ev.client.send)
				}
			}
		case <-ctx.Done():
			return
		}
	}
}

// Handler upgrades the HTTP connection to WebSocket and registers the client
// with the hub. Requires:
//   - auth.UserContextMiddleware to have run (provides userID)
//   - a :groupID route param, e.g. /ws/groups/:groupID
func (h *Hub) Handler(echoCtx echo.Context, ctx context.Context) error {
	uc := usercontext.Get(echoCtx.Request().Context())
	if !uc.IsLoggedIn {
		return echo.ErrUnauthorized
	}

	conn, err := upgrader.Upgrade(echoCtx.Response(), echoCtx.Request(), nil)
	if err != nil {
		return err
	}

	groupID := echoCtx.Param("groupID")

	client := newClient(conn, h, uc.ID, uc.Username, groupID)
	h.events <- clientEvent{kind: eventAdd, client: client}

	go client.ReadPump(echoCtx)
	go client.WritePump(echoCtx, ctx)

	return nil
}

// broadcast sends a message to every client currently viewing the same group.
// If a client's send buffer is full it is dropped to protect other clients.
func (h *Hub) broadcast(msg Message) error {
	for _, client := range h.clients {
		if client.GroupID != msg.GroupID {
			continue
		}
		select {
		case client.send <- msg:
		default:
			delete(h.clients, client.ID)
			close(client.send)
		}
	}
	return nil
}
