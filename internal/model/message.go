package model

// Message represents the core content of a message sent by a user.
// There is exactly one Message per "logical message".
type Message struct {
	Base
	Content  string `gorm:"type:text;not null"`
	SenderID string `gorm:"type:uuid;not null"`
	Sender   User   `gorm:"foreignKey:SenderID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	UserMessages []UserMessage `gorm:"foreignKey:MessageID"`
}

func init() {
	Register(&Message{})
}
