package handler

import (
	"net/http"

	"github.com/alextilot/golang-htmx-chatapp/model"
	"github.com/alextilot/golang-htmx-chatapp/services"
	"github.com/alextilot/golang-htmx-chatapp/web"
	"github.com/alextilot/golang-htmx-chatapp/web/components"
	"github.com/alextilot/golang-htmx-chatapp/web/views"
	"github.com/labstack/echo/v4"
)

func createMessageViews(messages []*model.Message, username string) []*components.MessageComponentViewModel {
	var messageViews []*components.MessageComponentViewModel
	for i := len(messages) - 1; i >= 0; i-- {
		msg := messages[i]
		tmp := components.NewMessageView(msg.Number, msg.Username, msg.Data, msg.Time, msg.Username == username)
		messageViews = append(messageViews, tmp)
	}
	return messageViews
}

func (h *Handler) Chatroom(etx echo.Context, user services.UserContext) error {
	messages, err := h.msgService.GetMostRecent(20)
	if err != nil {
		return etx.String(http.StatusBadGateway, "unable to pre populate chat messages")
	}
	messageViews := createMessageViews(messages, user.Username)
	return web.Render(etx, http.StatusOK, views.ChatroomPage(user.IsLoggedIn, messageViews))
}

func (h *Handler) Messages(etx echo.Context, user services.UserContext) error {
	req := &messagesRequest{}

	if err := req.bind(etx); err != nil {
		return etx.String(http.StatusBadRequest, err.Error())
	}

	messages, err := h.msgService.GetRows(req.Count, req.Ref)
	if err != nil {
		return etx.String(http.StatusInternalServerError, "unable to complete message request")
	}
	messageViews := createMessageViews(messages, user.Username)
	return web.Render(etx, http.StatusOK, components.MessageGroup(messageViews))
}
