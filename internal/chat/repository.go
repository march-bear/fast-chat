package chat

import (
	"errors"
	"fast-chat/internal/protocol"
	"sync"
	"time"
)

type Repository interface {
	addChatMessage(msg protocol.ChatMessage) (protocol.ChatMessage, error)
	getMessages(from int, to int) ([]protocol.ChatMessage, error)
}

func MakeRepository() Repository {
	return &repoImpl{Messages: make([]protocol.ChatMessage, 0, 16)}
}

type repoImpl struct {
	Mu        sync.Mutex
	IdCounter int
	Messages  []protocol.ChatMessage
}

func (r *repoImpl) addChatMessage(msg protocol.ChatMessage) (protocol.ChatMessage, error) {
	r.Mu.Lock()
	defer r.Mu.Unlock()

	msg.Id = r.IdCounter
	r.IdCounter++
	msg.CreatedAt = time.Now()

	r.Messages = append(r.Messages, msg)

	return msg, nil
}

func (r *repoImpl) getMessages(from int, to int) ([]protocol.ChatMessage, error) {
	if from < 0 || to < 0 || from > to {
		return nil, errors.New("`from` and `to` should be non-negative and from < to")
	}

	r.Mu.Lock()
	defer r.Mu.Unlock()

	to = min(to, len(r.Messages))
	messages := make([]protocol.ChatMessage, to-from)
	copy(messages, r.Messages[from:to])

	return messages, nil
}
