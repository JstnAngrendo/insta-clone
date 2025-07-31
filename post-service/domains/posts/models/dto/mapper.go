package dto

import (
	"github.com/jstnangrendo/instagram-clone/post-service/domains/posts/entities"
)

func ToPostDTO(e entities.Post) PostDTO {
	tags := make([]TagDTO, len(e.Tags))
	for i, t := range e.Tags {
		tags[i] = TagDTO{
			ID:   t.ID,
			Name: t.Name,
		}
	}

	return PostDTO{
		ID:           e.ID,
		UserID:       e.UserID,
		Caption:      e.Caption,
		ImageURL:     e.ImageURL,
		ThumbnailURL: e.ThumbnailURL,
		LikeCount:    e.LikeCount,
		CreatedAt:    e.CreatedAt,
		Tags:         tags,
	}
}
