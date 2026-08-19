package ws

import (
	"bytes"
	"context"
	"time"

	"github.com/alextilot/golang-htmx-chatapp/internal/realtime"
	"github.com/alextilot/golang-htmx-chatapp/web/components/chat"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v5"
)

const (
	pongWait     = 10 * time.Second
	pingInterval = 9 * time.Second

	messageTimeFormat = "2006-01-02 3:04:05 pm"

	// messageChannelBuffer: once full, the receiver is dropped rather than
	// blocking delivery to everyone else — see Hub.deliver.
	messageChannelBuffer = 32
)

type incomingMessage struct {
	Content string `json:"content"`
}

// Client represents a single connected WebSocket session.
//
// UserID ties this session to a model.User.ID so we can match it against
// realtime.Message.OwnerID when deciding whether to render a message as
// "self".
type Client struct {
	conn     *websocket.Conn
	hub      *realtime.Hub
	ID       string // unique session ID (not the user's DB ID)
	UserID   string // model.User.ID — set from auth session at connect time
	Username string // denormalised from the auth principal to stamp outgoing messages
	GroupID  string // model.Group.ID the client is currently viewing
	recv     <-chan realtime.Message
}

func newClient(conn *websocket.Conn, hub *realtime.Hub, userID string, username string, groupID string) *Client {
	id := uuid.New().String()
	return &Client{
		conn:     conn,
		hub:      hub,
		ID:       id,
		UserID:   userID,
		Username: username,
		GroupID:  groupID,
		recv:     hub.Register(id, groupID, messageChannelBuffer),
	}
}

// ReadPump takes ctx separately from echoCtx: net/http cancels a request's
// context as soon as its handler returns, which happens right after the
// upgrade, so echoCtx.Request().Context() would already be cancelled by the
// time this loop sends its first message.
func (c *Client) ReadPump(echoCtx *echo.Context, ctx context.Context) {
	defer func() {
		c.conn.Close()
		c.hub.Unregister(c.ID)
	}()

	if err := c.conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		echoCtx.Logger().Error(err.Error())
		return
	}

	c.conn.SetPongHandler(func(string) error {
		return c.conn.SetReadDeadline(time.Now().Add(pongWait))
	})

	for {
		var incoming incomingMessage
		if err := c.conn.ReadJSON(&incoming); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				echoCtx.Logger().Error(err.Error())
			}
			return
		}

		if incoming.Content == "" {
			continue
		}

		if err := c.hub.SendMessage(ctx, c.GroupID, c.UserID, c.Username, incoming.Content); err != nil {
			echoCtx.Logger().Error(err.Error())
			continue
		}
	}
}

// WritePump pumps outbound messages from the hub to the WebSocket connection.
// Each client runs WritePump in its own goroutine.
func (c *Client) WritePump(echoCtx *echo.Context, ctx context.Context) {
	defer func() {
		c.conn.Close()
		c.hub.Unregister(c.ID)
	}()

	ticker := time.NewTicker(pingInterval)
	defer ticker.Stop()

	for {
		select {
		case msg, ok := <-c.recv:
			if !ok {
				return
			}

			buf := &bytes.Buffer{}
			isSelf := msg.OwnerID == c.UserID
			chat.Message(msg.Username, msg.Data, msg.Time.Format(messageTimeFormat), isSelf).Render(ctx, buf)

			if err := c.conn.WriteMessage(websocket.TextMessage, buf.Bytes()); err != nil {
				echoCtx.Logger().Error(err.Error())
				return
			}

		case <-ticker.C:
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				echoCtx.Logger().Error(err.Error())
				return
			}

		case <-ctx.Done():
			return
		}
	}
}
