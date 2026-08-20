package repository

import (
	"context"
	"time"

	"github.com/alextilot/golang-htmx-chatapp/internal/model"
	"gorm.io/gorm"
)

// MessageRepository handles persistence for chat messages and their
// per-recipient delivery records (model.Message / model.UserMessage).
type MessageRepository struct {
	*BaseRepository[model.Message]
}

func NewMessageRepository(db *gorm.DB) *MessageRepository {
	return &MessageRepository{
		BaseRepository: NewBaseRepository[model.Message](db),
	}
}

// CreateWithRecipients persists a Message and one UserMessage delivery
// record per recipient, in a single transaction.
func (r *MessageRepository) CreateWithRecipients(
	ctx context.Context,
	msg *model.Message,
	groupID string,
	recipientIDs []string,
) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(msg).Error; err != nil {
			return err
		}

		for _, recipientID := range recipientIDs {
			um := &model.UserMessage{
				MessageID: msg.ID,
				OwnerID:   recipientID,
				GroupID:   &groupID,
			}
			if err := tx.Create(um).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// ListForGroup returns messages delivered to recipientID within groupID,
// most recent first. If before is non-zero, only messages strictly older
// than before are considered — used to page further back into history.
func (r *MessageRepository) ListForGroup(
	ctx context.Context,
	groupID string,
	recipientID string,
	before time.Time,
	limit int,
) ([]model.UserMessage, error) {
	var out []model.UserMessage

	q := r.db.WithContext(ctx).
		Preload("Message").
		Preload("Message.Sender").
		Where("group_id = ? AND owner_id = ?", groupID, recipientID)

	if !before.IsZero() {
		q = q.Where("created_at < ?", before)
	}

	err := q.Order("created_at DESC").Limit(limit).Find(&out).Error

	return out, err
}
