package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jstnangrendo/instagram-clone/timeline-service/handlers"
)

func SetupRouter(h *handlers.TimelineHandler) *gin.Engine {
	r := gin.Default()
	api := r.Group("/timeline")
	{
		api.GET("/", h.GetTimeline)
	}
	return r
}
