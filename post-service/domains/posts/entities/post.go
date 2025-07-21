package entities

import "time"

type Post struct {
	ID           uint      `json:"id"`
	UserID       uint      `json:"user_id"`
	Caption      string    `json:"caption"`
	ImageURL     string    `json:"image_url"`
	ThumbnailURL string    `json:"thumbnail_url"`
	LikeCount    int64     `json:"like_count"`
	CreatedAt    time.Time `json:"created_at"`
	Tags         []Tag     `json:"tags" gorm:"many2many:post_tags;"`
}
