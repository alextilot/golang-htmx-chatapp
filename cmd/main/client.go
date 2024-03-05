package main

import (
	"bytes"
	"context"
	"fmt"
	"golang-app/web/components"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type WebSocketMessage struct {
	Content string
}

type Message struct {
	ClientID string
	Name     string
	Time     time.Time
	Data     string
}

type Client struct {
	Conn           *websocket.Conn
	ID             string
	Chatroom       string
	Manager        *Manager
	MessageChannel chan Message
}

var (
	pongWaitTime = time.Second * 10
	pingInterval = time.Second * 9
)

func NewClient(conn *websocket.Conn, manager *Manager) *Client {
	return &Client{
		Conn:           conn,
		ID:             uuid.New().String(),
		Chatroom:       "general",
		Manager:        manager,
		MessageChannel: make(chan Message),
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

		fmt.Printf("%s\n", message)
		if err := c.Manager.WriteMessage(Message{
			ClientID: c.ID,
			Name:     "Alex",
			Time:     time.Now(),
			Data:     message.Content,
		}, "general"); err != nil {
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
			components.Message(msg.Name, msg.Data, msg.Time.Format("2006-01-02 3:4:5 pm"), msg.ClientID == c.ID).Render(ctx, buffer)
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
