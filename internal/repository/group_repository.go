package repository

import (
	"context"

	"github.com/alextilot/golang-htmx-chatapp/internal/model"
	"gorm.io/gorm"
)

type GroupRepository struct {
	*BaseRepository[model.Group]
}

func NewGroupRepository(db *gorm.DB) *GroupRepository {
	return &GroupRepository{NewBaseRepository[model.Group](db)}
}

// FindByMember returns all groups a user belongs to, ordered newest first.
func (r *GroupRepository) FindByMember(ctx context.Context, userID string) ([]model.Group, error) {
	var groups []model.Group
	err := r.Find(ctx, &groups, func(db *gorm.DB) *gorm.DB {
		return db.
			Joins("JOIN user_groups ON user_groups.group_id = groups.id").
			Where("user_groups.user_id = ?", userID).
			Where("groups.is_active = ?", true).
			Order("groups.created_at DESC")
	})
	return groups, err
}
