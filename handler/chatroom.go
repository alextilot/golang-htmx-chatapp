package handler

import (
	"net/http"

	"github.com/alextilot/golang-htmx-chatapp/services"
	"github.com/alextilot/golang-htmx-chatapp/web"
	"github.com/alextilot/golang-htmx-chatapp/web/components"
	"github.com/alextilot/golang-htmx-chatapp/web/views"
	"github.com/labstack/echo/v4"
)

func (h *Handler) Chatroom(etx echo.Context, user services.UserContext) error {
	messages, err := h.msgService.GetMostRecent(10)
	if err != nil {
		return etx.String(http.StatusBadGateway, "unable to pre populate chat messages")
	}

	var messageViews []*components.MessageComponentViewModel
	for i := len(messages) - 1; i >= 0; i-- {
		msg := messages[i]
		tmp := components.NewMessageView(msg.Username, msg.Data, msg.Time, msg.Username == user.Username)
		messageViews = append(messageViews, tmp)
	}
	return web.Render(etx, http.StatusOK, views.ChatroomPage(user.IsLoggedIn, messageViews))
}
