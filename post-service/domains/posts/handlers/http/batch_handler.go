package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jstnangrendo/instagram-clone/post-service/domains/posts/models/dto"
	"github.com/jstnangrendo/instagram-clone/post-service/domains/posts/usecases"
)

type BatchRequest struct {
	IDs []uint `json:"ids"`
}

func BatchGetHandler(pu usecases.PostUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req BatchRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
			return
		}
		posts := make([]dto.PostDTO, 0, len(req.IDs))
		for _, id := range req.IDs {
			post, err := pu.GetByID(id)
			if err != nil {
				continue
			}
			posts = append(posts, dto.ToPostDTO(*post))
		}
		c.JSON(http.StatusOK, posts)
	}
}
