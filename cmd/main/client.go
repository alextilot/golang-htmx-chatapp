package main

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/alextilot/golang-htmx-chatapp/model"
	"github.com/alextilot/golang-htmx-chatapp/web/components"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type websocketData struct {
	Type    string
	Version string
	Content string
}

type WsNotification struct {
	ClientID string
	Name     string
}

type Client struct {
	Conn                *websocket.Conn
	ID                  string
	Chatroom            string
	Manager             *Manager
	MessageChannel      chan model.Message
	NotificationChannel chan WsNotification
	Name                string
}

var (
	pongWaitTime = time.Second * 10
	pingInterval = time.Second * 9
)

func NewClient(conn *websocket.Conn, manager *Manager, name string) *Client {
	return &Client{
		Conn:                conn,
		ID:                  uuid.New().String(),
		Chatroom:            "general",
		Manager:             manager,
		MessageChannel:      make(chan model.Message),
		NotificationChannel: make(chan WsNotification),
		Name:                name,
	}
}

func (c *Client) handleWsNotification(ctx echo.Context, data websocketData) {
	fmt.Println("handleWsNotification()")
}

func (c *Client) handleWsMessage(ctx echo.Context, data websocketData) {
	fmt.Println("handleWsMessage()")
	wsMsg := data

	// err := json.Unmarshal(data.Body, &wsMsg)
	// if err != nil {
	// 	ctx.Logger().Error(err)
	// }

	s := strings.TrimSpace(wsMsg.Content)
	if s == "" {
		return
	}

	model := model.Message{
		Number:   1,
		Username: c.Name,
		Time:     time.Now(),
		Data:     wsMsg.Content,
	}

	// Save message to db
	dbModel, err := c.Manager.messageService.Create(&model)
	if err != nil {
		ctx.Logger().Error(err)
		return
	}

	// Send message to other people
	if err = c.Manager.WriteMessage(*dbModel, "general"); err != nil {
		ctx.Logger().Error(err)
		return
	}
}

func (c *Client) ReadHandler(ctx echo.Context) {
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
		var wsData websocketData

		err := c.Conn.ReadJSON(&wsData)
		fmt.Println("data", wsData)
		if err != nil {
			ctx.Logger().Error(err)
			return
		}

		fmt.Println(wsData)

		switch wsData.Type {
		case "notification":
			c.handleWsNotification(ctx, wsData)
		case "message":
			c.handleWsMessage(ctx, wsData)
		}
	}
}

func (c *Client) WriteHandler(echoContext echo.Context, ctx context.Context) {
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
		// case notif, ok := <-c.NotificationChannel:
		// 	if !ok {
		// 		return
		// 	}

		case msg, ok := <-c.MessageChannel:
			if !ok {
				return
			}

			buffer := &bytes.Buffer{}
			msgView := components.NewMessageView(
				msg.Number,
				msg.Username,
				msg.Data,
				msg.Time,
				msg.Username == c.Name,
			)
			components.Message(msgView).Render(ctx, buffer)

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
