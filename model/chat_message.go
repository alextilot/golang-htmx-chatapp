package model

type ChatMessage struct {
	Base
	UserID      string    `gorm:"type:uuid;primary_key;"`
	User        User      `gorm:"foreignKey:UserID"` // Belongs-to relationship with User
	ChatGroupID string    `gorm:"type:uuid;primary_key;"`
	ChatGroup   ChatGroup `gorm:"foreignKey:GroupID"` // Belings-to relationship with ChatGroup
	Content     string    `gorm:"type:text;not null"` // Non-nullable message content
}
