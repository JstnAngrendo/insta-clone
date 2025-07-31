package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jstnangrendo/instagram-clone/timeline-service/client"
	"github.com/jstnangrendo/instagram-clone/timeline-service/domains/timeline/usecases"
)

type TimelineHandler struct {
	uc      usecases.TimelineUseCase
	postCli *client.PostClient
}

func NewTimelineHandler(uc usecases.TimelineUseCase, pc *client.PostClient) *TimelineHandler {
	return &TimelineHandler{uc: uc, postCli: pc}
}

func (h *TimelineHandler) GetTimeline(c *gin.Context) {
	userIDStr := c.GetString("userID")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	ids, err := h.uc.FetchTimeline(c, userID, 20)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	posts, err := h.postCli.BatchGet(c, ids)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, posts)
}
