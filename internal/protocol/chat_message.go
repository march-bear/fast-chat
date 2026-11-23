package protocol

import "time"

type ChatMessage struct {
	Id        int       `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Author    string    `json:"author"`
	AuthorId  int       `json:"author_id"`
	Text      string    `json:"text"`
}
