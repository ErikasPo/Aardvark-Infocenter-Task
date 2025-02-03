package domain

import (
	"sync"
)

type Message struct {
	ID   int
	Data string
}

type Topic struct {
	Name     string
	clients  map[chan Message]bool
	messages []Message
	mutex    sync.Mutex
}

const maxCacheSize = 10

func NewTopic(name string) *Topic {
	return &Topic{
		Name:     name,
		clients:  make(map[chan Message]bool),
		messages: []Message{},
	}
}

func (t *Topic) AddMessage(msg Message) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if len(t.messages) >= maxCacheSize {
		t.messages = t.messages[1:]
	}
	t.messages = append(t.messages, msg)

	for client := range t.clients {
		select {
		case client <- msg:
		default:
			close(client)
			delete(t.clients, client)
		}
	}
}

func (t *Topic) Subscribe(client chan Message) {
	t.mutex.Lock()
	t.clients[client] = true
	for _, msg := range t.messages {
		select {
		case client <- msg:
		default:
		}
	}
	t.mutex.Unlock()
}

func (t *Topic) Unsubscribe(client chan Message) {
	t.mutex.Lock()
	delete(t.clients, client)
	close(client)
	t.mutex.Unlock()
}
