package model

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
