package entities

import "time"

type PostCreatedEvent struct {
	PostID    int64     `json:"post_id"`
	AuthorID  int64     `json:"user_id"`
	Caption   string    `json:"caption,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
