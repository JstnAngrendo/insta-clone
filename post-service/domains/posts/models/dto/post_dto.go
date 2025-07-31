package dto

import "time"

type TagDTO struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type PostDTO struct {
	ID           uint      `json:"id"`
	UserID       uint      `json:"user_id"`
	Caption      string    `json:"caption"`
	ImageURL     string    `json:"image_url"`
	ThumbnailURL string    `json:"thumbnail_url"`
	LikeCount    int64     `json:"like_count"`
	CreatedAt    time.Time `json:"created_at"`
	Tags         []TagDTO  `json:"tags"`
}
