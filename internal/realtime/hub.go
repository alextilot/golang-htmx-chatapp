package realtime

import (
	"context"

	"github.com/alextilot/golang-htmx-chatapp/internal/model"
	"github.com/alextilot/golang-htmx-chatapp/internal/service"
)

type eventKind int

const (
	eventRegister eventKind = iota
	eventUnregister
	eventBroadcast
)

type receiver struct {
	id      string
	groupID string
	send    chan Message
}

type event struct {
	kind     eventKind
	receiver *receiver
	id       string
	message  Message
}

// ChatSender is the subset of service.GroupService the hub needs to persist
// an incoming chat message before broadcasting it to connected receivers.
// Messages are always persisted first — a broadcast is a delivery
// notification for a message that already exists, never the system of
// record for it.
type ChatSender interface {
	SendMessage(ctx context.Context, groupID string, senderID string, input service.SendMessageInput) (*model.Message, error)
}

// Hub maintains the set of active receivers and routes messages between
// them. All state is owned by the single goroutine running Run — every
// mutation, including broadcast delivery, is serialized through the events
// channel, so nothing here needs a mutex.
type Hub struct {
	receivers map[string]*receiver
	events    chan event
	chat      ChatSender
}

func NewHub(chat ChatSender) *Hub {
	return &Hub{
		receivers: make(map[string]*receiver),
		events:    make(chan event),
		chat:      chat,
	}
}

// Run processes hub events until ctx is cancelled. Callers must start this
// in its own goroutine (go hub.Run(ctx)) before calling Register or
// SendMessage — events is unbuffered, so both block until Run is consuming.
func (h *Hub) Run(ctx context.Context) {
	for {
		select {
		case ev := <-h.events:
			switch ev.kind {
			case eventRegister:
				h.receivers[ev.receiver.id] = ev.receiver
			case eventUnregister:
				if r, ok := h.receivers[ev.id]; ok {
					delete(h.receivers, ev.id)
					close(r.send)
				}
			case eventBroadcast:
				h.deliver(ev.message)
			}
		case <-ctx.Done():
			return
		}
	}
}

func (h *Hub) Register(id, groupID string, bufferSize int) <-chan Message {
	r := &receiver{id: id, groupID: groupID, send: make(chan Message, bufferSize)}
	h.events <- event{kind: eventRegister, receiver: r}
	return r.send
}

// Unregister is safe to call even if the receiver was already dropped for
// being too slow.
func (h *Hub) Unregister(id string) {
	h.events <- event{kind: eventUnregister, id: id}
}

func (h *Hub) SendMessage(ctx context.Context, groupID string, senderID string, username string, content string) error {
	msg, err := h.chat.SendMessage(ctx, groupID, senderID, service.SendMessageInput{Content: content})
	if err != nil {
		return err
	}

	h.events <- event{kind: eventBroadcast, message: Message{
		MessageID: msg.ID,
		OwnerID:   senderID,
		Username:  username,
		GroupID:   groupID,
		Time:      msg.CreatedAt,
		Data:      msg.Content,
	}}
	return nil
}

// deliver sends msg to every receiver registered for msg.GroupID. Only ever
// called from Run's goroutine, so it's free to mutate h.receivers directly.
// A receiver whose buffer is full is dropped rather than allowed to block
// delivery to everyone else.
func (h *Hub) deliver(msg Message) {
	for id, r := range h.receivers {
		if r.groupID != msg.GroupID {
			continue
		}
		select {
		case r.send <- msg:
		default:
			delete(h.receivers, id)
			close(r.send)
		}
	}
}
