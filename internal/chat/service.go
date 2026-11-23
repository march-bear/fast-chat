package chat

import "fast-chat/internal/protocol"

type Service interface {
	GetHistory(from int, to int) ([]protocol.ChatMessage, error)
	AddMessage(text string, authorId int) (protocol.ChatMessage, error)
}

func MakeService() Service {
	return &serviceImpl{ChatMessagesRepo: MakeRepository()}
}

type serviceImpl struct {
	ChatMessagesRepo Repository
}

func (s *serviceImpl) GetHistory(from int, to int) ([]protocol.ChatMessage, error) {
	messages, err := s.ChatMessagesRepo.getMessages(from, to)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

func (s *serviceImpl) AddMessage(text string, authorId int) (protocol.ChatMessage, error) {
	msg := protocol.ChatMessage{
		Author:   "Anonymous",
		AuthorId: authorId,
		Text:     text,
	}

	msg, err := s.ChatMessagesRepo.addChatMessage(msg)

	return msg, err
}
