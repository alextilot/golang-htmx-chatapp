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

// Deactivate soft-deletes a group by marking it inactive, preserving its
// row (and message history) rather than removing it.
func (r *GroupRepository) Deactivate(ctx context.Context, id string) error {
	return r.DB().WithContext(ctx).
		Model(&model.Group{}).
		Where("id = ?", id).
		Update("is_active", false).Error
}

// UpdateDetails updates a group's name and/or description. An empty string
// leaves the corresponding field unchanged, so partial updates (e.g. name
// only) don't require the caller to re-fetch and re-supply every field.
func (r *GroupRepository) UpdateDetails(ctx context.Context, id string, name string, description string) error {
	updates := map[string]any{}
	if name != "" {
		updates["name"] = name
	}
	if description != "" {
		updates["description"] = description
	}
	if len(updates) == 0 {
		return nil
	}

	return r.DB().WithContext(ctx).
		Model(&model.Group{}).
		Where("id = ?", id).
		Updates(updates).Error
}
