package http

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jstnangrendo/instagram-clone/post-service/domains/posts/entities"
	"github.com/jstnangrendo/instagram-clone/post-service/domains/posts/models/responses"
	"github.com/jstnangrendo/instagram-clone/post-service/domains/posts/usecases"
	infrastructure "github.com/jstnangrendo/instagram-clone/post-service/infrastructure/storage"
	"github.com/jstnangrendo/instagram-clone/post-service/utils"
)

type PostHandler struct {
	uc usecases.PostUsecase
}

func NewPostHandler(uc usecases.PostUsecase) *PostHandler {
	return &PostHandler{uc: uc}
}

func (h *PostHandler) Create(c *gin.Context) {
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

	post, err := h.uc.CreateWithTags(userID, caption, imageURL, thumbPath, tagEntities)
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

func (h *PostHandler) GetByID(c *gin.Context) {
	param, _ := strconv.ParseUint(c.Param("post_id"), 10, 64)
	post, err := h.uc.GetByID(uint(param))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
		return
	}
	c.JSON(http.StatusOK, responses.PostResponse{
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

func (h *PostHandler) GetByUser(c *gin.Context) {
	raw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	logged := raw.(uint)

	param, err := strconv.ParseUint(c.Param("user_id"), 10, 64)
	if err != nil || uint(param) != logged {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	page, size := utils.GetPaginationParams(c)
	posts, total, err := h.uc.GetByUserPaginated(logged, page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := make([]responses.PostResponse, len(posts))
	for i, post := range posts {
		resp[i] = responses.PostResponse{
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
		TotalPages: int((total + int64(size) - 1) / int64(size)),
		TotalItems: total,
		Data:       resp,
	})
}

func (h *PostHandler) Delete(c *gin.Context) {
	raw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := raw.(uint)

	param, _ := strconv.ParseUint(c.Param("post_id"), 10, 64)
	er := h.uc.Delete(uint(param), userID)
	if er != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": er.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "post deleted"})
}

func (h *PostHandler) Timeline(c *gin.Context) {
	raw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := raw.(uint)
	userIDStr := fmt.Sprintf("%d", userID)

	posts, err := h.uc.GetTimeline(userIDStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get timeline"})
		return
	}

	page, size := utils.GetPaginationParams(c)
	total := len(posts)
	offset := (page - 1) * size
	end := offset + size
	if end > total {
		end = total
	}

	slice := posts[offset:end]
	resp := make([]responses.PostResponse, len(slice))
	for i, post := range slice {
		resp[i] = responses.PostResponse{
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
	totalPages := (total + size - 1) / size

	c.JSON(http.StatusOK, responses.PaginationResponse{
		Page:       page,
		Size:       size,
		TotalPages: totalPages,
		TotalItems: int64(total),
		Data:       resp,
	})
}

func (h *PostHandler) GetByTag(c *gin.Context) {
	tag := c.Param("tagName")
	page, size := utils.GetPaginationParams(c)
	posts, total, err := h.uc.GetPostsByTagPaginated(tag, page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve posts"})
		return
	}

	resp := make([]responses.PostResponse, len(posts))
	for i, post := range posts {
		resp[i] = responses.PostResponse{
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
		TotalPages: int((total + int64(size) - 1) / int64(size)),
		TotalItems: total,
		Data:       resp,
	})
}

func extractTagNames(tags []entities.Tag) []string {
	names := make([]string, len(tags))
	for i, t := range tags {
		names[i] = t.Name
	}
	return names
}
