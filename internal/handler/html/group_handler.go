package html

import (
	"context"
	"net/http"
	"time"

	"github.com/alextilot/golang-htmx-chatapp/internal/apperr"
	"github.com/alextilot/golang-htmx-chatapp/internal/auth"
	"github.com/alextilot/golang-htmx-chatapp/internal/model"
	"github.com/alextilot/golang-htmx-chatapp/internal/routes"
	"github.com/alextilot/golang-htmx-chatapp/internal/service"
	"github.com/alextilot/golang-htmx-chatapp/web/components/chat"
	"github.com/alextilot/golang-htmx-chatapp/web/pages"
	"github.com/labstack/echo/v5"
)

// historyPageSize is how many messages a group page (or a scroll-up page
// load) fetches at once. Getting back fewer than this means history's
// start has been reached — see GroupHandler.Show and .Messages.
const historyPageSize = 30

// GroupHandler handles HTML/HTMX requests for groups.
//
// WebSocket connections for group chat are owned by handler/ws — this
// handler only deals with the HTML surface (list/create/update/delete/
// membership), so it depends on the group service, not the WS hub.
type GroupHandler struct {
	svc groupServicer
}

// groupServicer is the subset of GroupService GroupHandler needs.
type groupServicer interface {
	ListForUser(ctx context.Context, userID string) ([]model.Group, error)
	GetByID(ctx context.Context, id string) (*model.Group, error)
	Create(ctx context.Context, creatorID string, input service.CreateGroupInput) (*model.Group, error)
	Update(ctx context.Context, groupID string, requesterID string, input service.UpdateGroupInput) (*model.Group, error)
	Delete(ctx context.Context, groupID string, requesterID string) error
	AddMember(ctx context.Context, groupID string, requesterID string, input service.AddMemberInput) error
	ListMessages(ctx context.Context, groupID string, userID string, before time.Time, limit int) ([]model.UserMessage, error)
}

func NewGroupHandler(svc groupServicer) *GroupHandler {
	return &GroupHandler{svc: svc}
}

// List renders the list of groups for the logged-in user.
// GET /groups
func (h *GroupHandler) List(c *echo.Context) error {
	p := auth.PrincipalFromEcho(c)

	groups, err := h.svc.ListForUser(c.Request().Context(), p.ID)
	if err != nil {
		return err // unexpected — let middleware handle
	}

	return sendResponse(c, Response{
		Status: http.StatusOK,
		Page:   pages.GroupListPage(groups, p.Username),
	})
}

// Show renders a single group's chat page.
// GET /groups/:groupID
func (h *GroupHandler) Show(c *echo.Context) error {
	ctx := c.Request().Context()
	p := auth.PrincipalFromEcho(c)
	groupID := c.Param("groupID")

	group, err := h.svc.GetByID(ctx, groupID)
	if err != nil {
		if appErr, ok := apperr.As(err); ok {
			return sendResponse(c, Response{Status: statusFor(appErr.Kind)})
		}
		return err // unexpected — let middleware handle
	}

	groups, err := h.svc.ListForUser(ctx, p.ID)
	if err != nil {
		return err // unexpected — let middleware handle
	}

	// ListMessages returns newest first; the page renders top-to-bottom
	// chronologically, so this is oldest-to-newest ahead of the live
	// messages the WebSocket appends after connecting.
	history, err := h.svc.ListMessages(ctx, groupID, p.ID, time.Time{}, historyPageSize)
	if err != nil {
		return err // unexpected — let middleware handle
	}
	hasMore := len(history) == historyPageSize
	for i, j := 0, len(history)-1; i < j; i, j = i+1, j-1 {
		history[i], history[j] = history[j], history[i]
	}

	return sendResponse(c, Response{
		Status: http.StatusOK,
		Page:   pages.GroupPage(*group, groups, history, hasMore, p.ID, p.Username),
	})
}

// Messages returns an older page of a group's messages as an HTMX fragment,
// for scroll-up pagination — older than the "before" query param (an
// RFC3339Nano timestamp).
// GET /groups/:groupID/messages
func (h *GroupHandler) Messages(c *echo.Context) error {
	ctx := c.Request().Context()
	p := auth.PrincipalFromEcho(c)
	groupID := c.Param("groupID")

	var before time.Time
	if raw := c.QueryParam("before"); raw != "" {
		t, err := time.Parse(time.RFC3339Nano, raw)
		if err != nil {
			return sendResponse(c, Response{Status: http.StatusBadRequest})
		}
		before = t
	}

	history, err := h.svc.ListMessages(ctx, groupID, p.ID, before, historyPageSize)
	if err != nil {
		return err // unexpected — let middleware handle
	}
	hasMore := len(history) == historyPageSize
	for i, j := 0, len(history)-1; i < j; i, j = i+1, j-1 {
		history[i], history[j] = history[j], history[i]
	}

	return sendResponse(c, Response{
		Status:   http.StatusOK,
		Fragment: chat.HistoryPage(groupID, history, p.ID, hasMore),
	})
}

// Create handles new group form submission.
// POST /groups
func (h *GroupHandler) Create(c *echo.Context) error {
	p := auth.PrincipalFromEcho(c)

	// 1. Bind
	req, err := BindInput[GroupCreateRequest](c)
	if err != nil {
		fe := apperr.FieldErrors{apperr.FieldGlobal: {"We couldn't read that request."}}
		return sendResponse(c, Response{Status: http.StatusBadRequest, Errors: fe})
	}

	// 2. Create group — sanitizing, validation, and business rules all
	// happen inside the service.
	if _, err := h.svc.Create(c.Request().Context(), p.ID, service.CreateGroupInput{
		Name: req.Name,
		Type: req.Type,
	}); err != nil {
		if appErr, ok := apperr.As(err); ok {
			fe := fieldsOf(appErr)
			return sendResponse(c, Response{Status: statusFor(appErr.Kind), Errors: fe})
		}
		return err // unexpected — let middleware handle
	}

	// 3. Return response (HTMX partial, or redirect)
	return sendResponse(c, Response{
		Status:   http.StatusOK,
		Redirect: routes.Routes.Groups.Path,
	})
}

// Update handles group rename / description change.
// PUT /groups/:groupID
func (h *GroupHandler) Update(c *echo.Context) error {
	p := auth.PrincipalFromEcho(c)

	// 1. Bind
	req, err := BindInput[GroupUpdateRequest](c)
	if err != nil {
		fe := apperr.FieldErrors{apperr.FieldGlobal: {"We couldn't read that request."}}
		return sendResponse(c, Response{Status: http.StatusBadRequest, Errors: fe})
	}

	// 2. Update group
	if _, err := h.svc.Update(c.Request().Context(), c.Param("groupID"), p.ID, service.UpdateGroupInput{
		Name:        req.Name,
		Description: req.Description,
	}); err != nil {
		if appErr, ok := apperr.As(err); ok {
			fe := fieldsOf(appErr)
			return sendResponse(c, Response{Status: statusFor(appErr.Kind), Errors: fe})
		}
		return err // unexpected — let middleware handle
	}

	return sendResponse(c, Response{
		Status:   http.StatusOK,
		Redirect: routes.Routes.Groups.Path,
	})
}

// Delete soft-deletes a group.
// DELETE /groups/:groupID
func (h *GroupHandler) Delete(c *echo.Context) error {
	p := auth.PrincipalFromEcho(c)

	if err := h.svc.Delete(c.Request().Context(), c.Param("groupID"), p.ID); err != nil {
		if appErr, ok := apperr.As(err); ok {
			return sendResponse(c, Response{Status: statusFor(appErr.Kind)})
		}
		return err // unexpected — let middleware handle
	}

	return sendResponse(c, Response{
		Status:   http.StatusOK,
		Redirect: routes.Routes.Groups.Path,
	})
}

// AddMember adds a user to a group.
// POST /groups/:groupID/members
func (h *GroupHandler) AddMember(c *echo.Context) error {
	p := auth.PrincipalFromEcho(c)

	// 1. Bind
	req, err := BindInput[GroupAddMemberRequest](c)
	if err != nil {
		fe := apperr.FieldErrors{apperr.FieldGlobal: {"We couldn't read that request."}}
		return sendResponse(c, Response{Status: http.StatusBadRequest, Errors: fe})
	}

	// 2. Add member
	if err := h.svc.AddMember(c.Request().Context(), c.Param("groupID"), p.ID, service.AddMemberInput{
		UserID: req.UserID,
	}); err != nil {
		if appErr, ok := apperr.As(err); ok {
			fe := fieldsOf(appErr)
			return sendResponse(c, Response{Status: statusFor(appErr.Kind), Errors: fe})
		}
		return err // unexpected — let middleware handle
	}

	return sendResponse(c, Response{
		Status:   http.StatusOK,
		Redirect: "/groups/" + c.Param("groupID"),
	})
}
