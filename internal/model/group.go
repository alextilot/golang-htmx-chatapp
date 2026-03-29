package model

// Group types
const (
	GroupTypeDirect = "direct" // 1:1 chat
	GroupTypeGroup  = "group"  // Multi-user group chat
)

// Group represents a chat room, either 1:1 or a larger group.
type Group struct {
	Base
	Name        string `gorm:"type:text;not null"`         // Required group name
	Description string `gorm:"type:text"`                  // Optional description
	Type        string `gorm:"type:text;default:'direct'"` // direct or group
	IsActive    bool   `gorm:"default:true"`               // Active vs archived

	UserGroups []UserGroup   `gorm:"foreignKey:GroupID"` // Participants
	Messages   []UserMessage `gorm:"foreignKey:GroupID"` // Messages
}

func init() {
	Register(&Group{})
}

