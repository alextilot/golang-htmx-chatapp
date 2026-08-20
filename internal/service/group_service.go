package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/alextilot/golang-htmx-chatapp/internal/apperr"
	"github.com/alextilot/golang-htmx-chatapp/internal/model"
	"github.com/alextilot/golang-htmx-chatapp/internal/repository"
)

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

// CreateGroupInput is GroupService.Create's own input shape.
type CreateGroupInput struct {
	Name string `validate:"required,min=1,max=100"`
	Type string `validate:"required,oneof=direct group"`
}

// UpdateGroupInput is GroupService.Update's own input shape.
type UpdateGroupInput struct {
	Name        string `validate:"omitempty,min=1,max=100"`
	Description string `validate:"omitempty,max=500"`
}

// AddMemberInput is GroupService.AddMember's own input shape.
type AddMemberInput struct {
	UserID string `validate:"required"`
}

// SendMessageInput is GroupService.SendMessage's own input shape.
type SendMessageInput struct {
	Content string `validate:"required,max=5000"`
}

// requireMember returns a Forbidden failure if userID is not a member of
// groupID. Centralizing this check keeps every membership-gated operation
// below consistent instead of each hand-rolling the same lookup.
func (s *GroupService) requireMember(ctx context.Context, groupID string, userID string) error {
	ok, err := s.userGroups.IsMember(ctx, groupID, userID)
	if err != nil {
		return err
	}
	if !ok {
		return apperr.Forbidden("you are not a member of this group")
	}
	return nil
}

// Create makes a new group and adds the creator as its first member.
// For direct (1:1) groups, pass GroupTypeDirect and both user IDs as members.
func (s *GroupService) Create(ctx context.Context, creatorID string, input CreateGroupInput) (*model.Group, error) {
	input.Name = strings.TrimSpace(input.Name)
	input.Type = strings.TrimSpace(strings.ToLower(input.Type))

	if appErr := apperr.ValidateStruct(input); appErr != nil {
		return nil, appErr
	}

	var created *model.Group

	err := s.groups.Transaction(ctx, func(tx *repository.BaseRepository[model.Group]) error {
		g := &model.Group{
			Name:     input.Name,
			Type:     input.Type,
			IsActive: true,
		}
		if err := tx.Create(ctx, g); err != nil {
			return err
		}
		created = g

		// Add the creator as the first member, on the same transaction —
		// s.userGroups holds the outer (non-transactional) connection, which
		// would deadlock against the write lock this transaction is still
		// holding on SQLite.
		return repository.NewUserGroupRepository(tx.DB()).AddMember(ctx, g.ID, creatorID)
	})

	return created, err
}

// GetByID returns a group by ID.
func (s *GroupService) GetByID(ctx context.Context, id string) (*model.Group, error) {
	var g model.Group
	if err := s.groups.FindByID(ctx, id, &g); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, apperr.NotFound("group not found")
		}
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
func (s *GroupService) Update(ctx context.Context, groupID string, requesterID string, input UpdateGroupInput) (*model.Group, error) {
	input.Name = strings.TrimSpace(input.Name)
	input.Description = strings.TrimSpace(input.Description)

	if appErr := apperr.ValidateStruct(input); appErr != nil {
		return nil, appErr
	}

	if err := s.requireMember(ctx, groupID, requesterID); err != nil {
		return nil, err
	}

	if err := s.groups.UpdateDetails(ctx, groupID, input.Name, input.Description); err != nil {
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
func (s *GroupService) AddMember(ctx context.Context, groupID string, requesterID string, input AddMemberInput) error {
	input.UserID = strings.TrimSpace(input.UserID)

	if appErr := apperr.ValidateStruct(input); appErr != nil {
		return appErr
	}

	if err := s.requireMember(ctx, groupID, requesterID); err != nil {
		return err
	}
	return s.userGroups.AddMember(ctx, groupID, input.UserID)
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
func (s *GroupService) SendMessage(ctx context.Context, groupID string, senderID string, input SendMessageInput) (*model.Message, error) {
	input.Content = strings.TrimSpace(input.Content)

	if appErr := apperr.ValidateStruct(input); appErr != nil {
		return nil, appErr
	}

	if err := s.requireMember(ctx, groupID, senderID); err != nil {
		return nil, err
	}

	recipientIDs, err := s.userGroups.ListMemberIDs(ctx, groupID)
	if err != nil {
		return nil, err
	}

	msg := &model.Message{
		Content:  input.Content,
		SenderID: senderID,
	}

	if err := s.messages.CreateWithRecipients(ctx, msg, groupID, recipientIDs); err != nil {
		return nil, err
	}

	return msg, nil
}

// ListMessages returns the messages delivered to userID in groupID, most
// recent first, older than before (zero value = no cutoff).
func (s *GroupService) ListMessages(ctx context.Context, groupID string, userID string, before time.Time, limit int) ([]model.UserMessage, error) {
	if limit <= 0 {
		limit = 50
	}
	return s.messages.ListForGroup(ctx, groupID, userID, before, limit)
}
