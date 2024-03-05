package main

import (
	"context"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type ClientListEvent struct {
	EventType string
	Client    *Client
}

type Manager struct {
	ClientList             []*Client
	ClientListEventChannel chan *ClientListEvent
}

var (
	upgrader = websocket.Upgrader{}
)

func NewManager() *Manager {
	return &Manager{
		ClientList:             []*Client{},
		ClientListEventChannel: make(chan *ClientListEvent),
	}
}

func (m *Manager) HandleClientListEventChannel(ctx context.Context) {
	for {
		select {
		case clientListEvent, ok := <-m.ClientListEventChannel:
			if !ok {
				return
			}
			switch clientListEvent.EventType {
			case "ADD":
				for _, client := range m.ClientList {
					if client.ID == clientListEvent.Client.ID {
						return
					}
				}
				m.ClientList = append(m.ClientList, clientListEvent.Client)
			case "REMOVE":
				newSlice := []*Client{}
				for _, client := range m.ClientList {
					if client.ID == clientListEvent.Client.ID {
						continue
					}
					newSlice = append(newSlice, client)
				}
				m.ClientList = newSlice
			}
		case <-ctx.Done():
			return
		}
	}
}

func (m *Manager) Handler(etx echo.Context, ctx context.Context) error {
	ws, err := upgrader.Upgrade(etx.Response(), etx.Request(), nil)
	if err != nil {
		return err
	}

	newClient := NewClient(ws, m)

	m.ClientListEventChannel <- &ClientListEvent{
		EventType: "ADD",
		Client:    newClient,
	}

	go newClient.ReadMessages(etx)
	go newClient.WriteMessage(etx, ctx)

	return nil
}

func (m *Manager) WriteMessage(message Message, chatroom string) error {
	for _, client := range m.ClientList {
		if client.Chatroom != chatroom {
			continue
		}
		client.MessageChannel <- message
	}
	return nil
}
