package model

// UserContact represents a contact/friend relationship between two users.
// Each row links a user to a contact user.
// The combination of UserID + ContactID is unique to prevent duplicate contacts.
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
