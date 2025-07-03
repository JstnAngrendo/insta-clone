package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jstnangrendo/instagram-clone/user-service/domains/users/entities"
	"github.com/jstnangrendo/instagram-clone/user-service/domains/users/models/requests"
	"github.com/jstnangrendo/instagram-clone/user-service/domains/users/models/responses"
	"github.com/jstnangrendo/instagram-clone/user-service/domains/users/usecases"
)

var userUC usecases.UserUsecase

func RegisterHandler(uc usecases.UserUsecase) gin.HandlerFunc {
	userUC = uc
	return func(c *gin.Context) {
		var req requests.RegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := userUC.Register(&entities.User{
			Username: req.Username,
			Email:    req.Email,
			Password: req.Password,
		}); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to register"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"message": "registered"})
	}
}

func LoginHandler(uc usecases.UserUsecase) gin.HandlerFunc {
	userUC = uc
	return func(c *gin.Context) {
		var req requests.LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		token, err := userUC.Login(req.Email, req.Password)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, responses.LoginResponse{Token: token})
	}
}

func ProfileHandler(uc usecases.UserUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		raw, _ := c.Get("user_id")
		userID := raw.(uint)

		user, err := uc.GetProfile(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load profile"})
			return
		}

		resp := responses.ProfileResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		}
		c.JSON(http.StatusOK, resp)
	}
}
