package model

// UserGroup is a pure join table between users and groups.
// Each row represents a unique membership.
type UserGroup struct {
	Base

	UserID string `gorm:"type:uuid;not null;primaryKey;index"`
	User   User   `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	GroupID string `gorm:"type:uuid;not null;primaryKey;index"`
	Group   Group  `gorm:"foreignKey:GroupID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func init() {
	Register(&UserGroup{})
}
