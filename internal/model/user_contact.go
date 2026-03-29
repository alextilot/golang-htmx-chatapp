package model

type UserContact struct {
	Base

	UserID string `gorm:"type:uuid;not null;index;uniqueIndex:idx_user_contact"`
	User   User   `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	ContactID string `gorm:"type:uuid;not null;index;uniqueIndex:idx_user_contact"`
	Contact   User   `gorm:"foreignKey:ContactID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func init() {
	Register(&UserContact{})
}
