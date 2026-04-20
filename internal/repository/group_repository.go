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

// IsMember returns true if the user belongs to the group.
func (r *GroupRepository) IsMember(ctx context.Context, groupID string, userID string) (bool, error) {
	var membership model.UserGroup
	err := r.First(ctx, &membership, func(db *gorm.DB) *gorm.DB {
		return db.Where("group_id = ? AND user_id = ?", groupID, userID)
	})
	if err != nil {
		return false, nil // not found = not a member
	}
	return true, nil
}

// AddMember inserts a UserGroup membership row.
func (r *GroupRepository) AddMember(ctx context.Context, groupID string, userID string) error {
	membership := &model.UserGroup{
		GroupID: groupID,
		UserID:  userID,
	}
	return r.db.WithContext(ctx).Create(membership).Error
}
