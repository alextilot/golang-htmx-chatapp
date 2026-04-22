package service

import (
	"context"
	"errors"

	"github.com/alextilot/golang-htmx-chatapp/internal/model"
	"github.com/alextilot/golang-htmx-chatapp/internal/repository"
)

// GroupService handles business logic for group management.
type GroupService struct {
	groups     *repository.GroupRepository
	userGroups *repository.UserGroupRepository
}

func NewGroupService(groups *repository.GroupRepository, userGroups *repository.UserGroupRepository) *GroupService {
	return &GroupService{groups: groups, userGroups: userGroups}
}

// Create makes a new group and adds the creator as its first member.
// For direct (1:1) groups, pass GroupTypeDirect and both user IDs as members.
func (s *GroupService) Create(ctx context.Context, creatorID string, name string, groupType string) (*model.Group, error) {
	var created *model.Group

	err := s.groups.Transaction(ctx, func(tx *repository.BaseRepository[model.Group]) error {
		g := &model.Group{
			Name:     name,
			Type:     groupType,
			IsActive: true,
		}
		if err := tx.Create(ctx, g); err != nil {
			return err
		}
		created = g

		// Add the creator as the first member.
		return s.userGroups.AddMember(ctx, g.ID, creatorID)
	})

	return created, err
}

// GetByID returns a group by ID. Returns an error if not found.
func (s *GroupService) GetByID(ctx context.Context, id string) (*model.Group, error) {
	var g model.Group
	if err := s.groups.FindByID(ctx, id, &g); err != nil {
		return nil, err
	}
	return &g, nil
}

// ListForUser returns all active groups the user is a member of.
func (s *GroupService) ListForUser(ctx context.Context, userID string) ([]model.Group, error) {
	return s.groups.FindByMember(ctx, userID)
}

// Update renames or changes the description of a group.
// Requires the requester to be a member.
// TODO: Add a role field to UserGroup to restrict this to admins/owners.
func (s *GroupService) Update(ctx context.Context, groupID string, requesterID string, name string, description string) (*model.Group, error) {
	ok, err := s.userGroups.IsMember(ctx, groupID, requesterID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("user is not a member of this group")
	}

	updates := map[string]any{}
	if name != "" {
		updates["name"] = name
	}
	if description != "" {
		updates["description"] = description
	}
	if len(updates) == 0 {
		return s.GetByID(ctx, groupID)
	}

	if err := s.groups.Update(ctx, groupID, updates); err != nil {
		return nil, err
	}
	return s.GetByID(ctx, groupID)
}

// Delete soft-deletes a group by marking it inactive.
// Requires the requester to be a member.
// TODO: Add a role field to UserGroup to restrict this to owners only.
func (s *GroupService) Delete(ctx context.Context, groupID string, requesterID string) error {
	ok, err := s.userGroups.IsMember(ctx, groupID, requesterID)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("user is not a member of this group")
	}
	// Soft-delete by marking inactive rather than removing the row,
	// so message history is preserved.
	return s.groups.Update(ctx, groupID, map[string]any{"is_active": false})
}

// AddMember adds a user to an existing group.
func (s *GroupService) AddMember(ctx context.Context, groupID string, requesterID string, newUserID string) error {
	ok, err := s.userGroups.IsMember(ctx, groupID, requesterID)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("user is not a member of this group")
	}
	return s.userGroups.AddMember(ctx, groupID, newUserID)
}
