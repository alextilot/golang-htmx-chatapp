package model

type UserMessage struct {
	Base
	Content string `gorm:"type:text;not null"` // Non-nullable message content

	OwnerID string `gorm:"type:uuid;not null"`
	Owner   User   `gorm:"foreignKey:OwnerID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	GroupID string `gorm:"type:uuid;not null"`
	Group   Group  `gorm:"foreignKey:GroupID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

// Notes:
// - UpdatedAt can be used as "edited at" since only Content is expected to change.
// - Expected that ID's won't change due to cascade.

func init() {
	Register(&UserMessage{})
}
