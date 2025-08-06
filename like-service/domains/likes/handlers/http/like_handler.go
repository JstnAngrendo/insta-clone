package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jstnangrendo/instagram-clone/like-service/domains/likes/usecases"
)

type LikeHandler struct {
	uc *usecases.LikeUseCase
}

func NewLikeHandler(uc *usecases.LikeUseCase) *LikeHandler {
	return &LikeHandler{uc: uc}
}

func (h *LikeHandler) LikePost(c *gin.Context) {
	userID := c.GetUint("user_id")
	pidStr := c.Param("post_id")
	pid, err := strconv.ParseUint(pidStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post_id"})
		return
	}
	err = h.uc.LikePost(c, uint(pid), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to like post"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "post liked"})
}

func (h *LikeHandler) UnlikePost(c *gin.Context) {
	userID := c.GetUint("user_id")
	pidStr := c.Param("post_id")
	pid, err := strconv.ParseUint(pidStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post_id"})
		return
	}
	err = h.uc.UnlikePost(c, uint(pid), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to unlike post"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "post unliked"})
}

func (h *LikeHandler) GetLikeCount(c *gin.Context) {
	pidStr := c.Param("post_id")
	pid, err := strconv.ParseUint(pidStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post_id"})
		return
	}
	count, err := h.uc.GetLikeCount(c, uint(pid))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get like count"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"like_count": count})
}

func (h *LikeHandler) GetPopularPosts(c *gin.Context) {
	limitStr := c.Query("limit")
	if limitStr == "" {
		limitStr = "10"
	}
	lim, err := strconv.ParseInt(limitStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit"})
		return
	}
	ids, err := h.uc.GetPopularPosts(c, lim)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get popular posts"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"popular_post_ids": ids})
}
