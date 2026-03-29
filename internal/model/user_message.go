package model

import "time"

// UserMessage represents a copy of a message for a specific user or group.
// Tracks per-recipient read status.
type UserMessage struct {
	Base
	MessageID string  `gorm:"type:uuid;not null"`
	Message   Message `gorm:"foreignKey:MessageID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	OwnerID string `gorm:"type:uuid;not null"`
	Owner   User   `gorm:"foreignKey:OwnerID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	GroupID *string `gorm:"type:uuid"` // Optional: nil for 1:1 DM
	Group   *Group  `gorm:"foreignKey:GroupID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	ReadAt *time.Time
}

func init() {
	Register(&UserMessage{})
}
