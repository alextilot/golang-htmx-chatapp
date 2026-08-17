package service

import (
	"context"
	"errors"

	"github.com/alextilot/golang-htmx-chatapp/internal/model"
	"github.com/alextilot/golang-htmx-chatapp/internal/repository"
)

// ErrNotMember is returned when an operation requires the actor to already
// be a member of the group, and they are not.
var ErrNotMember = errors.New("user is not a member of this group")

// GroupService handles business logic for groups, including chat within a
// group. A "chat room" is a Group — there is no separate chat entity in the
// data model, so a single service (used by both the html and ws handler
// packages) owns membership and messaging together, rather than two
// services independently re-implementing the same membership checks
// against the same table.
type GroupService struct {
	groups     *repository.GroupRepository
	userGroups *repository.UserGroupRepository
	messages   *repository.MessageRepository
}

func NewGroupService(
	groups *repository.GroupRepository,
	userGroups *repository.UserGroupRepository,
	messages *repository.MessageRepository,
) *GroupService {
	return &GroupService{groups: groups, userGroups: userGroups, messages: messages}
}

// requireMember returns ErrNotMember if userID is not a member of groupID.
// Centralizing this check keeps every membership-gated operation below
// consistent instead of each hand-rolling the same lookup and error.
func (s *GroupService) requireMember(ctx context.Context, groupID string, userID string) error {
	ok, err := s.userGroups.IsMember(ctx, groupID, userID)
	if err != nil {
		return err
	}
	if !ok {
		return ErrNotMember
	}
	return nil
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

// ListForUser returns all active groups (chat rooms) the user is a member of.
func (s *GroupService) ListForUser(ctx context.Context, userID string) ([]model.Group, error) {
	return s.groups.FindByMember(ctx, userID)
}

// Update renames or changes the description of a group.
// Requires the requester to be a member.
// TODO: Add a role field to UserGroup to restrict this to admins/owners.
func (s *GroupService) Update(ctx context.Context, groupID string, requesterID string, name string, description string) (*model.Group, error) {
	if err := s.requireMember(ctx, groupID, requesterID); err != nil {
		return nil, err
	}

	if err := s.groups.UpdateDetails(ctx, groupID, name, description); err != nil {
		return nil, err
	}
	return s.GetByID(ctx, groupID)
}

// Delete soft-deletes a group by marking it inactive.
// Requires the requester to be a member.
// TODO: Add a role field to UserGroup to restrict this to owners only.
func (s *GroupService) Delete(ctx context.Context, groupID string, requesterID string) error {
	if err := s.requireMember(ctx, groupID, requesterID); err != nil {
		return err
	}
	// Soft-delete by marking inactive rather than removing the row,
	// so message history is preserved.
	return s.groups.Deactivate(ctx, groupID)
}

// AddMember adds a user to an existing group. Requires the requester
// (the person doing the inviting) to already be a member.
func (s *GroupService) AddMember(ctx context.Context, groupID string, requesterID string, newUserID string) error {
	if err := s.requireMember(ctx, groupID, requesterID); err != nil {
		return err
	}
	return s.userGroups.AddMember(ctx, groupID, newUserID)
}

// Join adds userID to a group directly (self-service). Unlike AddMember,
// this does not require the actor to already be a member — that's the
// whole point — and is idempotent: joining a room you already belong to is
// not an error.
func (s *GroupService) Join(ctx context.Context, groupID string, userID string) error {
	ok, err := s.userGroups.IsMember(ctx, groupID, userID)
	if err != nil {
		return err
	}
	if ok {
		return nil
	}
	return s.userGroups.AddMember(ctx, groupID, userID)
}

// Leave removes userID from a group.
func (s *GroupService) Leave(ctx context.Context, groupID string, userID string) error {
	return s.userGroups.RemoveMember(ctx, groupID, userID)
}

// SendMessage persists a message and delivers it to every current member of
// the group, including the sender.
//
// Delivery is fan-out-on-write: one UserMessage row is written per
// recipient at send time (see MessageRepository.CreateWithRecipients and
// model.UserMessage). That's a deliberate, simple choice for small
// groups/DMs — it makes "list my messages" a single indexed lookup with no
// fan-in at read time. It does mean write cost scales linearly with group
// size, so if group sizes grow well beyond a handful of members, this
// should move to a read-time fan-out model instead (store the message
// once, resolve recipients at query time).
func (s *GroupService) SendMessage(ctx context.Context, groupID string, senderID string, content string) (*model.Message, error) {
	if content == "" {
		return nil, errors.New("message content is required")
	}

	if err := s.requireMember(ctx, groupID, senderID); err != nil {
		return nil, err
	}

	recipientIDs, err := s.userGroups.ListMemberIDs(ctx, groupID)
	if err != nil {
		return nil, err
	}

	msg := &model.Message{
		Content:  content,
		SenderID: senderID,
	}

	if err := s.messages.CreateWithRecipients(ctx, msg, groupID, recipientIDs); err != nil {
		return nil, err
	}

	return msg, nil
}

// ListMessages returns the messages delivered to userID in groupID, most
// recent first.
func (s *GroupService) ListMessages(ctx context.Context, groupID string, userID string, limit int) ([]model.UserMessage, error) {
	if limit <= 0 {
		limit = 50
	}
	return s.messages.ListForGroup(ctx, groupID, userID, limit)
}
