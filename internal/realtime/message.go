package realtime

import "time"

// Message is intentionally decoupled from the GORM model layer.
//
// Field mapping from model.UserMessage:
//
//	MessageID -> model.UserMessage.MessageID   (the persisted message's UUID)
//	OwnerID   -> model.UserMessage.OwnerID
//	Username  -> model.User.Username           (denormalised for display)
//	GroupID   -> model.UserMessage.GroupID
//	ReadAt    -> model.UserMessage.ReadAt      (nil = unread)
type Message struct {
	MessageID string
	OwnerID   string
	Username  string
	GroupID   string
	Time      time.Time
	Data      string
	ReadAt    *time.Time
}
