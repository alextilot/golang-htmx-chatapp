package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Base struct {
	// ID uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"` // For PostgreSQL
	// For other databases, remove default:uuid_generate_v4() and use BeforeCreate hook
	ID        string `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (b *Base) BeforeCreate(tx *gorm.DB) (err error) {
	b.ID = uuid.New().String() // Generates a new UUID and converts it to a string
	return
}
