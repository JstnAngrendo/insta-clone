package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jstnangrendo/instagram-clone/post-service/domains/posts/models/responses"
	"github.com/jstnangrendo/instagram-clone/post-service/domains/posts/usecases"
)

func CreatePostHandler(pu usecases.PostUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		raw, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		userID := raw.(uint)

		var req struct {
			Caption  string `json:"caption"`
			ImageURL string `json:"image_url"`
		}
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}

		post, err := pu.Create(userID, req.Caption, req.ImageURL)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, responses.PostResponse{
			ID:        post.ID,
			UserID:    post.UserID,
			Caption:   post.Caption,
			ImageURL:  post.ImageURL,
			CreatedAt: post.CreatedAt,
		})
	}
}

func GetPostHandler(pu usecases.PostUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		postID, _ := strconv.ParseUint(c.Param("post_id"), 10, 64)
		post, err := pu.GetByID(uint(postID))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
			return
		}
		c.JSON(http.StatusOK, post)
	}
}

func GetUserPostsHandler(pu usecases.PostUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDInterface, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found in context"})
			return
		}
		loggedInUserID := userIDInterface.(uint)

		paramUserID, err := strconv.ParseUint(c.Param("user_id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id parameter"})
			return
		}

		if uint(paramUserID) != loggedInUserID {
			c.JSON(http.StatusForbidden, gin.H{"error": "You are not allowed to access this resource"})
			return
		}

		posts, err := pu.GetByUser(loggedInUserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, posts)
	}
}

func DeletePostHandler(pu usecases.PostUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		raw, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		userID := raw.(uint)

		postID, _ := strconv.ParseUint(c.Param("post_id"), 10, 64)
		if err := pu.Delete(uint(postID), userID); err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "post deleted"})
	}
}

// Like & Unlike Post
func LikePostHandler(pu usecases.PostUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		raw, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		userID := raw.(uint)

		postID, _ := strconv.Atoi(c.Param("post_id"))

		err := pu.LikePost(userID, uint(postID))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "post liked"})
	}
}

func UnlikePostHandler(pu usecases.PostUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		raw, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		userID := raw.(uint)

		postID, _ := strconv.Atoi(c.Param("post_id"))

		err := pu.UnlikePost(userID, uint(postID))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "post unliked"})
	}
}

func GetPostLikesHandler(pu usecases.PostUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		postID, err := strconv.Atoi(c.Param("post_id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
			return
		}

		count, err := pu.CountLikes(uint(postID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to count likes"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"like_count": count})
	}
}
