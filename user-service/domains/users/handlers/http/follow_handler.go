package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jstnangrendo/instagram-clone/user-service/domains/users/usecases"
)

func FollowHandler(fu usecases.FollowUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		raw, _ := c.Get("user_id")
		userID := raw.(uint)
		targetID, err := strconv.ParseUint(c.Param("user_id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid target id"})
			return
		}
		if err := fu.Follow(userID, uint(targetID)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "followed"})
	}
}

func UnfollowHandler(fu usecases.FollowUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		raw, _ := c.Get("user_id")
		userID := raw.(uint)
		targetID, _ := strconv.ParseUint(c.Param("user_id"), 10, 64)
		if err := fu.Unfollow(userID, uint(targetID)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "unfollowed"})
	}
}

func GetFollowingIDsHandler(fu usecases.FollowUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}

		followingIDs, err := fu.GetFollowingUserIDs(uint(userID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"following_ids": followingIDs})
	}
}
