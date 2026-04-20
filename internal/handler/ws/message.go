package ws

import "time"

// incomingMessage is the raw JSON payload received from the client over WebSocket.
type incomingMessage struct {
	Content string `json:"content"`
}

// Message is the WS-layer DTO for a chat message routed through the hub.
//
// It is intentionally decoupled from the GORM model layer. Populate it from
// a model.UserMessage at the boundary (e.g. in your Echo handler or service
// layer) before passing it into the hub.
//
// Field mapping from model.UserMessage:
//
//	MessageID  -> model.UserMessage.MessageID   (the persisted message's UUID)
//	OwnerID    -> model.UserMessage.OwnerID     (replaces ClientID/Username)
//	Username   -> model.User.Username           (denormalised for display)
//	GroupID    -> model.UserMessage.GroupID     (replaces the string Chatroom)
//	SenderSelf -> compared at render time: msg.OwnerID == client.UserID
//	ReadAt     -> model.UserMessage.ReadAt      (nil = unread)
type Message struct {
	MessageID string
	OwnerID   string
	Username  string
	GroupID   string
	Time      time.Time
	Data      string
	ReadAt    *time.Time
}
