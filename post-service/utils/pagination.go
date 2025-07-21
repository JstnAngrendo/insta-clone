package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetPaginationParams(c *gin.Context) (page, size int) {
	pageStr := c.DefaultQuery("page", "1")
	sizeStr := c.DefaultQuery("size", "10")

	page, _ = strconv.Atoi(pageStr)
	size, _ = strconv.Atoi(sizeStr)

	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 10
	}
	return
}
