package handler

import (
	"context"
	"net/http"

	"github.com/alextilot/golang-htmx-chatapp/internal/handler/ws"
	"github.com/alextilot/golang-htmx-chatapp/internal/model"
	"github.com/alextilot/golang-htmx-chatapp/internal/response"
	"github.com/alextilot/golang-htmx-chatapp/internal/usercontext"
	"github.com/alextilot/golang-htmx-chatapp/internal/validation"
	"github.com/alextilot/golang-htmx-chatapp/internal/validation/schema"
	"github.com/labstack/echo/v4"
)

// GroupHandler handles HTTP and WebSocket requests for groups.
type GroupHandler struct {
	hub *ws.Hub
	svc groupServicer
}

// groupServicer is the subset of GroupService GroupHandler needs.
type groupServicer interface {
	ListForUser(ctx context.Context, userID string) ([]model.Group, error)
	GetByID(ctx context.Context, id string) (*model.Group, error)
	Create(ctx context.Context, creatorID string, name string, groupType string) (*model.Group, error)
	Update(ctx context.Context, groupID string, requesterID string, name string, description string) (*model.Group, error)
	Delete(ctx context.Context, groupID string, requesterID string) error
	AddMember(ctx context.Context, groupID string, requesterID string, newUserID string) error
}

func NewGroupHandler(hub *ws.Hub, svc groupServicer) *GroupHandler {
	return &GroupHandler{hub: hub, svc: svc}
}

// List renders the list of groups for the logged-in user.
// GET /app/groups
func (h *GroupHandler) List(c echo.Context) error {
	uc := usercontext.FromEcho(c)

	groups, err := h.svc.ListForUser(c.Request().Context(), uc.ID)
	if err != nil {
		fe := validation.FieldErrors{validation.FieldServer: {"Failed to load groups"}}
		return response.Send(c, response.Response{
			Status: http.StatusInternalServerError,
			Errors: fe,
		})
	}

	return response.Send(c, response.Response{
		Status: http.StatusOK,
		Data:   map[string]any{"groups": groups},
		// TODO: View: pages.GroupListPage(groups),
	})
}

// Show renders a single group's chat page.
// GET /app/groups/:groupID
func (h *GroupHandler) Show(c echo.Context) error {
	g, err := h.svc.GetByID(c.Request().Context(), c.Param("groupID"))
	if err != nil {
		fe := validation.FieldErrors{validation.FieldRequest: {"Group not found"}}
		return response.Send(c, response.Response{
			Status: http.StatusNotFound,
			Errors: fe,
		})
	}

	return response.Send(c, response.Response{
		Status: http.StatusOK,
		Data:   map[string]any{"group": g},
		// TODO: View: pages.GroupPage(g),
	})
}

// Create handles new group form submission.
// POST /app/groups
func (h *GroupHandler) Create(c echo.Context) error {
	uc := usercontext.FromEcho(c)

	// 1. Bind, sanitize, validate
	input, err := schema.HandleInput(c, schema.SanitizeGroupCreate, schema.ValidateGroupCreate)
	if err != nil {
		if ierr, ok := err.(*schema.InputError); ok {
			return response.Send(c, response.Response{
				Status: http.StatusBadRequest,
				// TODO: View: forms.GroupCreateForm(ierr.Fields),
				Errors: ierr.Fields,
			})
		}
		return err
	}

	// 2. Create group
	g, err := h.svc.Create(c.Request().Context(), uc.ID, input.Name, input.Type)
	if err != nil {
		fe := validation.FieldErrors{validation.FieldServer: {"Failed to create group"}}
		return response.Send(c, response.Response{
			Status: http.StatusInternalServerError,
			// TODO: View: forms.GroupCreateForm(fe),
			Errors: fe,
		})
	}

	// 3. Return response (HTMX partial, JSON, or redirect)
	return response.Send(c, response.Response{
		Status:   http.StatusOK,
		Data:     map[string]any{"group": g},
		Redirect: "/app/groups",
	})
}

// Update handles group rename / description change.
// PUT /app/groups/:groupID
func (h *GroupHandler) Update(c echo.Context) error {
	uc := usercontext.FromEcho(c)

	// 1. Bind, sanitize, validate
	input, err := schema.HandleInput(c, schema.SanitizeGroupUpdate, schema.ValidateGroupUpdate)
	if err != nil {
		if ierr, ok := err.(*schema.InputError); ok {
			return response.Send(c, response.Response{
				Status: http.StatusBadRequest,
				// TODO: View: forms.GroupEditForm(ierr.Fields),
				Errors: ierr.Fields,
			})
		}
		return err
	}

	// 2. Update group
	g, err := h.svc.Update(
		c.Request().Context(),
		c.Param("groupID"),
		uc.ID,
		input.Name,
		input.Description,
	)
	if err != nil {
		fe := validation.FieldErrors{validation.FieldServer: {"Failed to update group"}}
		return response.Send(c, response.Response{
			Status: http.StatusInternalServerError,
			// TODO: View: forms.GroupEditForm(fe),
			Errors: fe,
		})
	}

	return response.Send(c, response.Response{
		Status:   http.StatusOK,
		Data:     map[string]any{"group": g},
		Redirect: "/app/groups",
	})
}

// Delete soft-deletes a group.
// DELETE /app/groups/:groupID
func (h *GroupHandler) Delete(c echo.Context) error {
	uc := usercontext.FromEcho(c)

	if err := h.svc.Delete(c.Request().Context(), c.Param("groupID"), uc.ID); err != nil {
		fe := validation.FieldErrors{validation.FieldServer: {"Failed to delete group"}}
		return response.Send(c, response.Response{
			Status: http.StatusInternalServerError,
			Errors: fe,
		})
	}

	return response.Send(c, response.Response{
		Status:   http.StatusOK,
		Data:     map[string]any{"deleted": true},
		Redirect: "/app/groups",
	})
}

// AddMember adds a user to a group.
// POST /app/groups/:groupID/members
func (h *GroupHandler) AddMember(c echo.Context) error {
	uc := usercontext.FromEcho(c)

	// 1. Bind, sanitize, validate
	input, err := schema.HandleInput(c, schema.SanitizeGroupAddMember, schema.ValidateGroupAddMember)
	if err != nil {
		if ierr, ok := err.(*schema.InputError); ok {
			return response.Send(c, response.Response{
				Status: http.StatusBadRequest,
				Errors: ierr.Fields,
			})
		}
		return err
	}

	// 2. Add member
	if err := h.svc.AddMember(c.Request().Context(), c.Param("groupID"), uc.ID, input.UserID); err != nil {
		fe := validation.FieldErrors{validation.FieldServer: {"Failed to add member"}}
		return response.Send(c, response.Response{
			Status: http.StatusInternalServerError,
			Errors: fe,
		})
	}

	return response.Send(c, response.Response{
		Status:   http.StatusOK,
		Data:     map[string]any{"added": true},
		Redirect: "/app/groups/" + c.Param("groupID"),
	})
}

// ConnectWS upgrades the connection to WebSocket for a group's chat.
// GET /ws/groups/:groupID
func (h *GroupHandler) ConnectWS(c echo.Context, ctx context.Context) error {
	return h.hub.Handler(c, ctx)
}
