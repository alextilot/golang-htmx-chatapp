package main

import (
	"bytes"
	"context"
	"strings"
	"time"

	"github.com/alextilot/golang-htmx-chatapp/model"
	"github.com/alextilot/golang-htmx-chatapp/web/components"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type WebSocketMessage struct {
	Content string
}

type Client struct {
	Conn           *websocket.Conn
	ID             string
	Chatroom       string
	Manager        *Manager
	MessageChannel chan model.Message
	Name           string
}

var (
	pongWaitTime = time.Second * 10
	pingInterval = time.Second * 9
)

func NewClient(conn *websocket.Conn, manager *Manager, name string) *Client {
	return &Client{
		Conn:           conn,
		ID:             uuid.New().String(),
		Chatroom:       "general",
		Manager:        manager,
		MessageChannel: make(chan model.Message),
		Name:           name,
	}
}

func (c *Client) ReadMessages(ctx echo.Context) {
	defer func() {
		c.Conn.Close()
		c.Manager.ClientListEventChannel <- &ClientListEvent{
			Client:    c,
			EventType: "REMOVE",
		}
	}()

	if err := c.Conn.SetReadDeadline(time.Now().Add(pongWaitTime)); err != nil {
		ctx.Logger().Error(err)
		return
	}

	c.Conn.SetPongHandler(func(appData string) error {
		if err := c.Conn.SetReadDeadline(time.Now().Add(pongWaitTime)); err != nil {
			ctx.Logger().Error(err)
			return err
		}
		return nil
	})

	for {
		var message WebSocketMessage
		err := c.Conn.ReadJSON(&message)
		if err != nil {
			ctx.Logger().Error(err)
			return
		}

		s := strings.TrimSpace(message.Content)
		if s == "" {
			return
		}

		msg := model.Message{
			Number:   0,
			Username: c.Name,
			Time:     time.Now(),
			Data:     message.Content,
		}

		// Save message to db
		dbMsg, err := c.Manager.messageService.Create(&msg)
		if err != nil {
			ctx.Logger().Error(err)
			return
		}

		// Send message to other people
		if err = c.Manager.WriteMessage(*dbMsg, "general"); err != nil {
			ctx.Logger().Error(err)
			return
		}
	}
}

func (c *Client) WriteMessage(echoContext echo.Context, ctx context.Context) {
	defer func() {
		c.Conn.Close()
		c.Manager.ClientListEventChannel <- &ClientListEvent{
			Client:    c,
			EventType: "REMOVE",
		}
	}()

	//check for alive (heartbeat)
	ticker := time.NewTicker(pingInterval)

	for {
		select {
		case msg, ok := <-c.MessageChannel:
			if !ok {
				return
			}

			buffer := &bytes.Buffer{}
			tmp := components.NewMessageView(
				msg.Number,
				msg.Username,
				msg.Data,
				msg.Time,
				msg.Username == c.Name,
			)
			components.Message(tmp).Render(ctx, buffer)

			err := c.Conn.WriteMessage(websocket.TextMessage, buffer.Bytes())
			if err != nil {
				echoContext.Logger().Error(err)
				return
			}
		case <-ticker.C:
			if err := c.Conn.WriteMessage(websocket.PingMessage, []byte("")); err != nil {
				echoContext.Logger().Error(err)
				return
			}
		case <-ctx.Done():
			return
		}
	}
}
