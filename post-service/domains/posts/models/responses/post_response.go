package responses

import "time"

type PostResponse struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	Caption   string    `json:"caption"`
	ImageURL  string    `json:"image_url"`
	CreatedAt time.Time `json:"created_at"`
	LikeCount int       `json:"like_count"`
}
