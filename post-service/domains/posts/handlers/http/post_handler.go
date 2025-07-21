package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jstnangrendo/instagram-clone/post-service/domains/posts/entities"
	"github.com/jstnangrendo/instagram-clone/post-service/domains/posts/models/responses"
	"github.com/jstnangrendo/instagram-clone/post-service/domains/posts/usecases"
	infrastructure "github.com/jstnangrendo/instagram-clone/post-service/infrastructure/storage"
	"github.com/jstnangrendo/instagram-clone/post-service/utils"
)

func CreatePostHandler(pu usecases.PostUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		raw, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		userID := raw.(uint)

		caption := c.PostForm("caption")
		file, err := c.FormFile("image")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "image is required"})
			return
		}

		imageURL, err := infrastructure.UploadImage(c, file)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upload image"})
			return
		}

		thumbPath, err := utils.GenerateThumbnail("uploads/" + file.Filename)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate thumbnail"})
			return
		}

		formTags := c.PostFormArray("tags")
		tagEntities := utils.CleanTags(formTags)

		post, err := pu.CreateWithTags(userID, caption, imageURL, thumbPath, tagEntities)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, responses.PostResponse{
			ID:           post.ID,
			UserID:       post.UserID,
			Caption:      post.Caption,
			ImageURL:     post.ImageURL,
			ThumbnailURL: post.ThumbnailURL,
			CreatedAt:    post.CreatedAt,
			LikeCount:    post.LikeCount,
			Tags:         extractTagNames(post.Tags),
		})
	}
}

func extractTagNames(tags []entities.Tag) []string {
	names := make([]string, len(tags))
	for i, tag := range tags {
		names[i] = tag.Name
	}
	return names
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
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		loggedInUserID := userIDInterface.(uint)

		paramUserID, err := strconv.ParseUint(c.Param("user_id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id parameter"})
			return
		}
		if uint(paramUserID) != loggedInUserID {
			c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
			return
		}

		page, size := utils.GetPaginationParams(c)

		posts, total, err := pu.GetByUserPaginated(loggedInUserID, page, size)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		postResponses := make([]responses.PostResponse, len(posts))
		for i, post := range posts {
			postResponses[i] = responses.PostResponse{
				ID:           post.ID,
				UserID:       post.UserID,
				Caption:      post.Caption,
				ImageURL:     post.ImageURL,
				ThumbnailURL: post.ThumbnailURL,
				CreatedAt:    post.CreatedAt,
				LikeCount:    post.LikeCount,
				Tags:         extractTagNames(post.Tags),
			}
		}

		totalPages := int((total + int64(size) - 1) / int64(size))

		c.JSON(http.StatusOK, responses.PaginationResponse{
			Page:       page,
			Size:       size,
			TotalPages: totalPages,
			TotalItems: total,
			Data:       postResponses,
		})
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

func GetPostsByTagHandler(pu usecases.PostUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		tagName := c.Param("tagName")
		page, size := utils.GetPaginationParams(c)

		posts, total, err := pu.GetPostsByTagPaginated(tagName, page, size)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve posts"})
			return
		}
		totalPages := int((total + int64(size) - 1) / int64(size))

		postResponses := make([]responses.PostResponse, len(posts))
		for i, post := range posts {
			postResponses[i] = responses.PostResponse{
				ID:           post.ID,
				UserID:       post.UserID,
				Caption:      post.Caption,
				ImageURL:     post.ImageURL,
				ThumbnailURL: post.ThumbnailURL,
				CreatedAt:    post.CreatedAt,
				LikeCount:    post.LikeCount,
				Tags:         extractTagNames(post.Tags),
			}
		}

		c.JSON(http.StatusOK, responses.PaginationResponse{
			Page:       page,
			Size:       size,
			TotalPages: totalPages,
			TotalItems: total,
			Data:       postResponses,
		})
	}
}
