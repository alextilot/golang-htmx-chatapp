package repository

import (
	"context"
	"errors"

	"github.com/alextilot/golang-htmx-chatapp/internal/model"
	"gorm.io/gorm"
)

type UserGroupRepository struct {
	*BaseRepository[model.UserGroup]
}

func NewUserGroupRepository(db *gorm.DB) *UserGroupRepository {
	return &UserGroupRepository{NewBaseRepository[model.UserGroup](db)}
}

// IsMember only treats a genuine not-found as "not a member" — any other
// error (e.g. a cancelled context) is propagated rather than misread as one.
func (r *UserGroupRepository) IsMember(ctx context.Context, groupID string, userID string) (bool, error) {
	var membership model.UserGroup
	err := r.First(ctx, &membership, func(db *gorm.DB) *gorm.DB {
		return db.Where("group_id = ? AND user_id = ?", groupID, userID)
	})
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if err != nil {
		return false, err
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

// RemoveMember deletes a membership row.
func (r *UserGroupRepository) RemoveMember(ctx context.Context, groupID string, userID string) error {
	return r.db.WithContext(ctx).
		Where("group_id = ? AND user_id = ?", groupID, userID).
		Delete(&model.UserGroup{}).Error
}

// ListMemberIDs returns the user IDs of every member of a group.
func (r *UserGroupRepository) ListMemberIDs(ctx context.Context, groupID string) ([]string, error) {
	var ids []string
	err := r.db.WithContext(ctx).
		Model(&model.UserGroup{}).
		Where("group_id = ?", groupID).
		Pluck("user_id", &ids).Error
	return ids, err
}
