package model

type ChatGroup struct {
	Base
	Name        string `gorm:"type:text;not null"`
	Description *string
}
