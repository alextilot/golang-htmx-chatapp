package handler

import (
	"context"
	"net/http"

	"github.com/alextilot/golang-htmx-chatapp/internal/handler/ws"
	"github.com/alextilot/golang-htmx-chatapp/internal/model"
	"github.com/alextilot/golang-htmx-chatapp/internal/usercontext"
	"github.com/labstack/echo/v4"
)

// GroupHandler handles HTTP and WebSocket requests for groups.
type GroupHandler struct {
	hub *ws.Hub
	svc groupServicer
}

// groupServicer is the subset of GroupService GroupHandler needs.
// Using an interface keeps the handler testable without a real DB.
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
		return err
	}
	_ = groups
	// TODO: return web.Render(c, http.StatusOK, pages.GroupListPage(groups))
	return c.String(http.StatusOK, "group list")
}

// Show renders a single group's chat page.
// GET /app/groups/:groupID
func (h *GroupHandler) Show(c echo.Context) error {
	g, err := h.svc.GetByID(c.Request().Context(), c.Param("groupID"))
	if err != nil {
		return echo.ErrNotFound
	}
	_ = g
	// TODO: return web.Render(c, http.StatusOK, pages.GroupPage(g))
	return c.String(http.StatusOK, "group chat page")
}

// Create handles new group form submission.
// POST /app/groups
func (h *GroupHandler) Create(c echo.Context) error {
	uc := usercontext.FromEcho(c)

	name := c.FormValue("name")
	if name == "" {
		return echo.ErrBadRequest
	}

	groupType := c.FormValue("type")
	if groupType != model.GroupTypeDirect && groupType != model.GroupTypeGroup {
		groupType = model.GroupTypeGroup // default to multi-user
	}

	g, err := h.svc.Create(c.Request().Context(), uc.ID, name, groupType)
	if err != nil {
		return err
	}
	_ = g
	// TODO: return HTMX partial or redirect
	return c.Redirect(http.StatusSeeOther, "/app/groups")
}

// Update handles group rename / description change.
// PUT /app/groups/:groupID
func (h *GroupHandler) Update(c echo.Context) error {
	uc := usercontext.FromEcho(c)

	_, err := h.svc.Update(
		c.Request().Context(),
		c.Param("groupID"),
		uc.ID,
		c.FormValue("name"),
		c.FormValue("description"),
	)
	if err != nil {
		return err
	}
	return c.Redirect(http.StatusSeeOther, "/app/groups")
}

// Delete soft-deletes a group (marks it inactive).
// DELETE /app/groups/:groupID
func (h *GroupHandler) Delete(c echo.Context) error {
	uc := usercontext.FromEcho(c)

	if err := h.svc.Delete(c.Request().Context(), c.Param("groupID"), uc.ID); err != nil {
		return err
	}
	// TODO: return HTMX swap to remove the row from the list
	return c.Redirect(http.StatusSeeOther, "/app/groups")
}

// AddMember adds a user to a group.
// POST /app/groups/:groupID/members
func (h *GroupHandler) AddMember(c echo.Context) error {
	uc := usercontext.FromEcho(c)

	newUserID := c.FormValue("userID")
	if newUserID == "" {
		return echo.ErrBadRequest
	}

	if err := h.svc.AddMember(c.Request().Context(), c.Param("groupID"), uc.ID, newUserID); err != nil {
		return err
	}
	return c.Redirect(http.StatusSeeOther, "/app/groups/"+c.Param("groupID"))
}

// ConnectWS upgrades the connection to WebSocket for a group's chat.
// GET /ws/groups/:groupID
func (h *GroupHandler) ConnectWS(c echo.Context, ctx context.Context) error {
	return h.hub.Handler(c, ctx)
}
