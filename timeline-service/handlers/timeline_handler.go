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
	uid := c.Query("user_id")
	userID, err := strconv.ParseInt(uid, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}
	ids, err := h.uc.FetchTimeline(c, userID, 20)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	posts, _ := h.postCli.BatchGet(c, ids)
	c.JSON(http.StatusOK, posts)
}
