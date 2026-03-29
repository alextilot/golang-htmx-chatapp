package model

// Group types
const (
	GroupTypeDirect = "direct" // 1:1 chat
	GroupTypeGroup  = "group"  // Multi-user group chat
)

// Group represents a chat room, either a 1:1 conversation or a multi-user group.
type Group struct {
	Base
	Name        string `gorm:"type:text;not null"`         // Required group name
	Description string `gorm:"type:text"`                  // Optional description
	Type        string `gorm:"type:text;default:'direct'"` // 'direct' or 'group'
	IsActive    bool   `gorm:"default:true"`               // Active vs archived

	// Relationships
	UserGroups   []UserGroup   `gorm:"foreignKey:GroupID"` // Users in this group
	UserMessages []UserMessage `gorm:"foreignKey:GroupID"` // Messages sent to this group
}

func init() {
	Register(&Group{})
}

