package entities

import "time"

type Comment struct {
	ID        uint `gorm:"primaryKey" json:"id"`
	PostID    uint  `json:"post_id"`
	UserID    uint  `json:"user_id"`
	Content   string  `json:"content"`
	CreatedAt time.Time  `json:"created_at"`
}
