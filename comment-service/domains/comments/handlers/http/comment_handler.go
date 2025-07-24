package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jstnangrendo/instagram-clone/post-service/domains/comments/usecases"
)

func CreateCommentHandler(pu usecases.CommentUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		raw, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		userID := raw.(uint)

		postIDParam := c.Param("post_id")
		postID, err := strconv.ParseUint(postIDParam, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post_id"})
			return
		}

		var body struct {
			Content string `json:"content"`
		}
		if err := c.ShouldBindJSON(&body); err != nil || body.Content == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "content is required"})
			return
		}

		err = pu.CreateComment(uint(postID), userID, body.Content)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add comment"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "comment created"})
	}
}

func GetCommentsHandler(pu usecases.CommentUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		postIDParam := c.Param("post_id")
		postID, err := strconv.ParseUint(postIDParam, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post_id"})
			return
		}

		comments, err := pu.GetComments(uint(postID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch comments"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"comments": comments})
	}
}
