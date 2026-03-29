package model

// UserGroup is the join table connecting users to groups.
// Each row represents a membership of a user in a group.
// The combination of UserID + GroupID is unique to prevent duplicate memberships.
type UserGroup struct {
	Base

	UserID string `gorm:"type:uuid;not null;index;uniqueIndex:idx_user_group"`
	User   User   `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	GroupID string `gorm:"type:uuid;not null;index;uniqueIndex:idx_user_group"`
	Group   Group  `gorm:"foreignKey:GroupID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func init() {
	Register(&UserGroup{})
}
