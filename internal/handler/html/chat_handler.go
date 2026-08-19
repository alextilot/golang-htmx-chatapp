package html

import (
	"context"
	"net/http"

	"github.com/alextilot/golang-htmx-chatapp/internal/apperr"
	"github.com/alextilot/golang-htmx-chatapp/internal/auth"
	"github.com/alextilot/golang-htmx-chatapp/internal/model"
	"github.com/alextilot/golang-htmx-chatapp/internal/routes"
	"github.com/alextilot/golang-htmx-chatapp/internal/service"
	"github.com/alextilot/golang-htmx-chatapp/web/pages"
	"github.com/labstack/echo/v5"
)

// ChatHandler handles HTML/HTMX requests for chat rooms and messages.
//
// A chat room is a Group (see service.GroupService) — WebSocket delivery of
// new messages is owned by handler/ws, so this handler covers the HTML
// surface only (listing rooms, rendering the room page, and the
// join/leave/send actions).
type ChatHandler struct {
	svc chatServicer
}

// chatServicer is the subset of GroupService ChatHandler needs. Chat rooms
// are Groups (see service.GroupService) — this handler only touches the
// membership/messaging methods, not group management (create/update/
// delete), which is GroupHandler's concern.
type chatServicer interface {
	ListForUser(ctx context.Context, userID string) ([]model.Group, error)
	GetByID(ctx context.Context, id string) (*model.Group, error)
	Join(ctx context.Context, groupID string, userID string) error
	Leave(ctx context.Context, groupID string, userID string) error
	SendMessage(ctx context.Context, groupID string, senderID string, input service.SendMessageInput) (*model.Message, error)
}

func NewChatHandler(svc chatServicer) *ChatHandler {
	return &ChatHandler{svc: svc}
}

// List renders the chat rooms the logged-in user belongs to.
// GET /chat
func (h *ChatHandler) List(c *echo.Context) error {
	p := auth.PrincipalFromEcho(c)

	if _, err := h.svc.ListForUser(c.Request().Context(), p.ID); err != nil {
		return err // unexpected — let middleware handle
	}

	return sendResponse(c, Response{
		Status: http.StatusOK,
		// TODO: Page: pages.ChatRoomListPage(rooms),
	})
}

// Room renders a single chat room.
// GET /chat/:id
func (h *ChatHandler) Room(c *echo.Context) error {
	if _, err := h.svc.GetByID(c.Request().Context(), c.Param("id")); err != nil {
		if appErr, ok := apperr.As(err); ok {
			return sendResponse(c, Response{Status: statusFor(appErr.Kind)})
		}
		return err // unexpected — let middleware handle
	}

	return sendResponse(c, Response{
		Status: http.StatusOK,
		Page:   pages.ChatroomPage(),
	})
}

// Join adds the logged-in user to a chat room.
// POST /chat/:id/join
func (h *ChatHandler) Join(c *echo.Context) error {
	p := auth.PrincipalFromEcho(c)

	if err := h.svc.Join(c.Request().Context(), c.Param("id"), p.ID); err != nil {
		return err // unexpected — let middleware handle
	}

	return sendResponse(c, Response{
		Status:   http.StatusOK,
		Redirect: "/chat/" + c.Param("id"),
	})
}

// Leave removes the logged-in user from a chat room.
// POST /chat/:id/leave
func (h *ChatHandler) Leave(c *echo.Context) error {
	p := auth.PrincipalFromEcho(c)

	if err := h.svc.Leave(c.Request().Context(), c.Param("id"), p.ID); err != nil {
		return err // unexpected — let middleware handle
	}

	return sendResponse(c, Response{
		Status:   http.StatusOK,
		Redirect: routes.Routes.Chat.Path,
	})
}

// SendMessage posts a new message to a chat room.
//
// Delivery to connected WebSocket clients is handled separately by
// handler/ws — this endpoint covers the non-WS (HTMX/plain HTTP) send path
// and persistence.
// POST /chat/:id/message
func (h *ChatHandler) SendMessage(c *echo.Context) error {
	p := auth.PrincipalFromEcho(c)

	req, err := BindInput[ChatMessageRequest](c)
	if err != nil {
		fe := apperr.FieldErrors{apperr.FieldGlobal: {"We couldn't read that request."}}
		return sendResponse(c, Response{Status: http.StatusBadRequest, Errors: fe})
	}

	if _, err := h.svc.SendMessage(c.Request().Context(), c.Param("id"), p.ID, service.SendMessageInput{
		Content: req.Content,
	}); err != nil {
		if appErr, ok := apperr.As(err); ok {
			fe := fieldsOf(appErr)
			return sendResponse(c, Response{Status: statusFor(appErr.Kind), Errors: fe})
		}
		return err // unexpected — let middleware handle
	}

	return sendResponse(c, Response{
		Status: http.StatusOK,
	})
}
