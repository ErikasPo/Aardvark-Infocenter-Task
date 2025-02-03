package application

import (
	"Infocenter/Domain"
	"errors"
	"sync"
)

type MessageService struct {
	topics map[string]*domain.Topic
	mutex  sync.Mutex
	msgID  int
}

var errorEmptyMessage = errors.New("Please fill in the message.")

func NewMessageService() *MessageService {
	return &MessageService{
		topics: make(map[string]*domain.Topic),
	}
}

func (s *MessageService) GetTopic(name string) *domain.Topic {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if _, exists := s.topics[name]; !exists {
		s.topics[name] = domain.NewTopic(name)
	}
	return s.topics[name]
}

func (s *MessageService) PublishMessage(topicName, messageData string) (int, error) {
	if messageData == "" {
		return 0, errorEmptyMessage
	}

	s.mutex.Lock()
	s.msgID++
	message := domain.Message{ID: s.msgID, Data: messageData}
	s.mutex.Unlock()

	topic := s.GetTopic(topicName)
	topic.AddMessage(message)

	return message.ID, nil
}
