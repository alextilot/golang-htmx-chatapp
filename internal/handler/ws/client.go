package ws

import (
	"bytes"
	"context"
	"time"

	"github.com/alextilot/golang-htmx-chatapp/web/components"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

const (
	pongWait     = 10 * time.Second
	pingInterval = 9 * time.Second

	// messageTimeFormat is the display format for chat message timestamps.
	messageTimeFormat = "2006-01-02 3:04:05 pm"

	// messageChannelBuffer is the number of messages buffered per client
	// before the client is considered too slow and is dropped.
	messageChannelBuffer = 32
)

// Client represents a single connected WebSocket session.
//
// UserID ties this session to a model.User.ID so we can match it against
// Message.OwnerID when deciding whether to render a message as "self".
type Client struct {
	conn     *websocket.Conn
	hub      *Hub
	ID       string // unique session ID (not the user's DB ID)
	UserID   string // model.User.ID — set from auth session at connect time
	Username string // denormalised from usercontext to stamp outgoing messages
	GroupID  string // model.Group.ID the client is currently viewing
	send     chan Message
}

func newClient(conn *websocket.Conn, hub *Hub, userID string, username string, groupID string) *Client {
	return &Client{
		conn:     conn,
		hub:      hub,
		ID:       uuid.New().String(),
		UserID:   userID,
		Username: username,
		GroupID:  groupID,
		send:     make(chan Message, messageChannelBuffer),
	}
}

// ReadPump pumps inbound messages from the WebSocket connection to the hub.
// Each client runs ReadPump in its own goroutine.
func (c *Client) ReadPump(ctx echo.Context) {
	defer func() {
		c.conn.Close()
		c.hub.events <- clientEvent{client: c, kind: eventRemove}
	}()

	if err := c.conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		ctx.Logger().Error(err)
		return
	}

	c.conn.SetPongHandler(func(string) error {
		return c.conn.SetReadDeadline(time.Now().Add(pongWait))
	})

	for {
		var incoming incomingMessage
		if err := c.conn.ReadJSON(&incoming); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				ctx.Logger().Error(err)
			}
			return
		}

		if incoming.Content == "" {
			continue
		}

		// TODO: persist to DB here (insert model.Message + model.UserMessage),
		// then broadcast the returned MessageID so other clients can reference it.
		msg := Message{
			MessageID: uuid.New().String(), // replace with DB-assigned ID after persistence
			OwnerID:   c.UserID,
			Username:  c.Username,
			GroupID:   c.GroupID,
			Time:      time.Now(),
			Data:      incoming.Content,
		}
		if err := c.hub.broadcast(msg); err != nil {
			ctx.Logger().Error(err)
			return
		}
	}
}

// WritePump pumps outbound messages from the hub to the WebSocket connection.
// Each client runs WritePump in its own goroutine.
func (c *Client) WritePump(echoCtx echo.Context, ctx context.Context) {
	defer func() {
		c.conn.Close()
		c.hub.events <- clientEvent{client: c, kind: eventRemove}
	}()

	ticker := time.NewTicker(pingInterval)
	defer ticker.Stop()

	for {
		select {
		case msg, ok := <-c.send:
			if !ok {
				return
			}

			buf := &bytes.Buffer{}
			isSelf := msg.OwnerID == c.UserID
			components.Message(msg.Username, msg.Data, msg.Time.Format(messageTimeFormat), isSelf).Render(ctx, buf)

			if err := c.conn.WriteMessage(websocket.TextMessage, buf.Bytes()); err != nil {
				echoCtx.Logger().Error(err)
				return
			}

		case <-ticker.C:
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				echoCtx.Logger().Error(err)
				return
			}

		case <-ctx.Done():
			return
		}
	}
}
