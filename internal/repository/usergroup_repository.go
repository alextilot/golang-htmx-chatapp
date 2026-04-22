package repository

import (
	"context"

	"github.com/alextilot/golang-htmx-chatapp/internal/model"
	"gorm.io/gorm"
)

type UserGroupRepository struct {
	*BaseRepository[model.UserGroup]
}

func NewUserGroupRepository(db *gorm.DB) *UserGroupRepository {
	return &UserGroupRepository{NewBaseRepository[model.UserGroup](db)}
}

// IsMember returns true if the user belongs to the group.
func (r *UserGroupRepository) IsMember(ctx context.Context, groupID string, userID string) (bool, error) {
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
func (r *UserGroupRepository) AddMember(ctx context.Context, groupID string, userID string) error {
	membership := &model.UserGroup{
		GroupID: groupID,
		UserID:  userID,
	}
	return r.db.WithContext(ctx).Create(membership).Error
}
