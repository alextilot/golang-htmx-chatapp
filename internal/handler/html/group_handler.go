package html

import (
	"context"
	"net/http"

	"github.com/alextilot/golang-htmx-chatapp/internal/auth"
	"github.com/alextilot/golang-htmx-chatapp/internal/model"
	"github.com/alextilot/golang-htmx-chatapp/internal/routes"
	"github.com/alextilot/golang-htmx-chatapp/internal/validation"
	"github.com/alextilot/golang-htmx-chatapp/internal/validation/schema"
	"github.com/labstack/echo/v5"
)

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
	Create(ctx context.Context, creatorID string, name string, groupType string) (*model.Group, error)
	Update(ctx context.Context, groupID string, requesterID string, name string, description string) (*model.Group, error)
	Delete(ctx context.Context, groupID string, requesterID string) error
	AddMember(ctx context.Context, groupID string, requesterID string, newUserID string) error
}

func NewGroupHandler(svc groupServicer) *GroupHandler {
	return &GroupHandler{svc: svc}
}

// List renders the list of groups for the logged-in user.
// GET /groups
func (h *GroupHandler) List(c *echo.Context) error {
	p := auth.PrincipalFromEcho(c)

	_, err := h.svc.ListForUser(c.Request().Context(), p.ID)
	if err != nil {
		fe := validation.FieldErrors{validation.FieldServer: {"Failed to load groups"}}
		return sendResponse(c, Response{
			Status: http.StatusInternalServerError,
			Errors: fe,
		})
	}

	return sendResponse(c, Response{
		Status: http.StatusOK,
		// TODO: Page: pages.GroupListPage(groups),
	})
}

// Show renders a single group's chat page.
// GET /groups/:groupID
func (h *GroupHandler) Show(c *echo.Context) error {
	_, err := h.svc.GetByID(c.Request().Context(), c.Param("groupID"))
	if err != nil {
		fe := validation.FieldErrors{validation.FieldRequest: {"Group not found"}}
		return sendResponse(c, Response{
			Status: http.StatusNotFound,
			Errors: fe,
		})
	}

	return sendResponse(c, Response{
		Status: http.StatusOK,
		// TODO: Page: pages.GroupPage(g),
	})
}

// Create handles new group form submission.
// POST /groups
func (h *GroupHandler) Create(c *echo.Context) error {
	p := auth.PrincipalFromEcho(c)

	// 1. Bind, sanitize, validate
	input, err := schema.HandleInput(c, schema.SanitizeGroupCreate, schema.ValidateGroupCreate)
	if err != nil {
		if ierr, ok := err.(*schema.InputError); ok {
			return sendResponse(c, Response{
				Status: http.StatusBadRequest,
				// TODO: Page: forms.GroupCreateForm(ierr.Fields),
				Errors: ierr.Fields,
			})
		}
		return err
	}

	// 2. Create group
	if _, err := h.svc.Create(c.Request().Context(), p.ID, input.Name, input.Type); err != nil {
		fe := validation.FieldErrors{validation.FieldServer: {"Failed to create group"}}
		return sendResponse(c, Response{
			Status: http.StatusInternalServerError,
			// TODO: Page: forms.GroupCreateForm(fe),
			Errors: fe,
		})
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

	// 1. Bind, sanitize, validate
	input, err := schema.HandleInput(c, schema.SanitizeGroupUpdate, schema.ValidateGroupUpdate)
	if err != nil {
		if ierr, ok := err.(*schema.InputError); ok {
			return sendResponse(c, Response{
				Status: http.StatusBadRequest,
				// TODO: Page: forms.GroupEditForm(ierr.Fields),
				Errors: ierr.Fields,
			})
		}
		return err
	}

	// 2. Update group
	if _, err := h.svc.Update(
		c.Request().Context(),
		c.Param("groupID"),
		p.ID,
		input.Name,
		input.Description,
	); err != nil {
		fe := validation.FieldErrors{validation.FieldServer: {"Failed to update group"}}
		return sendResponse(c, Response{
			Status: http.StatusInternalServerError,
			// TODO: Page: forms.GroupEditForm(fe),
			Errors: fe,
		})
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
		fe := validation.FieldErrors{validation.FieldServer: {"Failed to delete group"}}
		return sendResponse(c, Response{
			Status: http.StatusInternalServerError,
			Errors: fe,
		})
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

	// 1. Bind, sanitize, validate
	input, err := schema.HandleInput(c, schema.SanitizeGroupAddMember, schema.ValidateGroupAddMember)
	if err != nil {
		if ierr, ok := err.(*schema.InputError); ok {
			return sendResponse(c, Response{
				Status: http.StatusBadRequest,
				Errors: ierr.Fields,
			})
		}
		return err
	}

	// 2. Add member
	if err := h.svc.AddMember(c.Request().Context(), c.Param("groupID"), p.ID, input.UserID); err != nil {
		fe := validation.FieldErrors{validation.FieldServer: {"Failed to add member"}}
		return sendResponse(c, Response{
			Status: http.StatusInternalServerError,
			Errors: fe,
		})
	}

	return sendResponse(c, Response{
		Status:   http.StatusOK,
		Redirect: "/groups/" + c.Param("groupID"),
	})
}
