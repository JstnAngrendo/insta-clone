package entities

import "time"

type PostCreatedEvent struct {
	PostID    string    `json:"post_id"`
	AuthorID  int64     `json:"author_id"`
	Timestamp time.Time `json:"timestamp"`
}
