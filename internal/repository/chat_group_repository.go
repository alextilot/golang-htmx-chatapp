package repository

import (
	"github.com/alextilot/golang-htmx-chatapp/internal/model"
	"gorm.io/gorm"
)

type ChatGroupRepository struct {
	*BaseRepository[model.ChatGroup]
}

func NewChatGroupRepository(db *gorm.DB) *ChatGroupRepository {
	return &ChatGroupRepository{
		BaseRepository: NewBaseRepository[model.ChatGroup](db),
	}
}

// Optional: type-specific query example
func (r *ChatGroupRepository) FindByName(name string) ([]model.ChatGroup, error) {
	var groups []model.ChatGroup
	if err := r.db.Where("name LIKE ?", "%"+name+"%").Find(&groups).Error; err != nil {
		return nil, err
	}
	return groups, nil
}
