package repository

import (
	"github.com/alextilot/golang-htmx-chatapp/internal/model"
	"gorm.io/gorm"
)

type ChatMessageRepository struct {
	*BaseRepository[model.ChatMessage]
}

func NewChatMessageRepository(db *gorm.DB) *ChatMessageRepository {
	return &ChatMessageRepository{
		BaseRepository: NewBaseRepository[model.ChatMessage](db),
	}
}

// FindByChatGroup retrieves all messages for a specific chat group
// Supports optional limit, offset (pagination), and ordering
func (r *ChatMessageRepository) FindByChatGroup(groupID string, limit, offset int, order string) ([]model.ChatMessage, error) {
	tx := r.db.Model(&model.ChatMessage{}).Where("chat_group_id = ?", groupID)

	if order != "" {
		tx = tx.Order(order)
	} else {
		tx = tx.Order("created_at asc") // default: oldest first
	}

	if limit > 0 {
		tx = tx.Limit(limit).Offset(offset)
	}

	var messages []model.ChatMessage
	if err := tx.Preload("User").Preload("ChatGroup").Find(&messages).Error; err != nil {
		return nil, err
	}

	return messages, nil
}

// FindByUser retrieves all messages sent by a specific user
func (r *ChatMessageRepository) FindByUser(userID string, limit, offset int, order string) ([]model.ChatMessage, error) {
	tx := r.db.Model(&model.ChatMessage{}).Where("user_id = ?", userID)

	if order != "" {
		tx = tx.Order(order)
	} else {
		tx = tx.Order("created_at asc")
	}

	if limit > 0 {
		tx = tx.Limit(limit).Offset(offset)
	}

	var messages []model.ChatMessage
	if err := tx.Preload("User").Preload("ChatGroup").Find(&messages).Error; err != nil {
		return nil, err
	}

	return messages, nil
}
