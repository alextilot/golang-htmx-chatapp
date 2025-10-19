package model

type ChatMessage struct {
	Base
	UserID      string    `gorm:"type:uuid"`
	User        User      `gorm:"foreignKey:UserID"` // Belongs-to relationship with User
	ChatGroupID string    `gorm:"type:uuid"`
	ChatGroup   ChatGroup `gorm:"foreignKey:ChatGroupID"` // Belings-to relationship with ChatGroup
	Content     string    `gorm:"type:text;not null"`     // Non-nullable message content
}

func init() {
	Register(&ChatMessage{})
}
