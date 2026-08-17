package ws

import (
	"context"
	"net/http"

	"github.com/alextilot/golang-htmx-chatapp/internal/auth"
	"github.com/alextilot/golang-htmx-chatapp/internal/model"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v5"
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

// ChatSender is the subset of service.GroupService the hub needs to
// persist a chat message and learn who it was delivered to, before
// broadcasting it to connected clients. Messages are always persisted
// first — a WebSocket broadcast is a delivery notification for a message
// that already exists, never the system of record for it.
type ChatSender interface {
	SendMessage(ctx context.Context, groupID string, senderID string, content string) (*model.Message, error)
}

// Hub maintains the set of active clients and routes messages between them.
// All client-list mutations are serialized through the events channel to
// avoid the need for a mutex.
type Hub struct {
	clients map[string]*Client
	events  chan clientEvent
	chat    ChatSender
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// TODO: restrict to your actual origin(s) before going to production.
		return true
	},
}

func NewHub(chat ChatSender) *Hub {
	return &Hub{
		clients: make(map[string]*Client),
		events:  make(chan clientEvent),
		chat:    chat,
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
// with the hub. Requires auth.Service.AuthMiddleware to have run. groupID identifies
// which group/room the connection should be scoped to; callers read it from
// whichever route param their endpoint uses (e.g. :groupID or :id).
func (h *Hub) Handler(echoCtx *echo.Context, ctx context.Context, groupID string) error {
	p := auth.PrincipalFromEcho(echoCtx)
	if !p.Authenticated {
		return echo.ErrUnauthorized
	}

	conn, err := upgrader.Upgrade(echoCtx.Response(), echoCtx.Request(), nil)
	if err != nil {
		return err
	}

	client := newClient(conn, h, p.ID, p.Username, groupID)
	h.events <- clientEvent{kind: eventAdd, client: client}

	go client.ReadPump(echoCtx)
	go client.WritePump(echoCtx, ctx)

	return nil
}

// SendMessage persists an incoming chat message via the chat service and
// broadcasts the persisted result to every client currently viewing the
// same group. This is the only path by which a WebSocket-originated
// message reaches other clients — there is no direct client-to-client
// broadcast that skips persistence.
func (h *Hub) SendMessage(ctx context.Context, groupID string, senderID string, username string, content string) error {
	msg, err := h.chat.SendMessage(ctx, groupID, senderID, content)
	if err != nil {
		return err
	}

	return h.broadcast(Message{
		MessageID: msg.ID,
		OwnerID:   senderID,
		Username:  username,
		GroupID:   groupID,
		Time:      msg.CreatedAt,
		Data:      msg.Content,
	})
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
