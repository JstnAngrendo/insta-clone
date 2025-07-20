package main

import (
	"os"

	"github.com/gin-gonic/gin"

	"github.com/jstnangrendo/instagram-clone/user-service/config"
	"github.com/jstnangrendo/instagram-clone/user-service/domains/users/entities"
	httpUser "github.com/jstnangrendo/instagram-clone/user-service/domains/users/handlers/http"
	"github.com/jstnangrendo/instagram-clone/user-service/domains/users/repositories"
	"github.com/jstnangrendo/instagram-clone/user-service/domains/users/usecases"
	"github.com/jstnangrendo/instagram-clone/user-service/middleware"
)

func main() {

	config.InitDatabase()
	config.InitRedis()

	config.DB.AutoMigrate(&entities.User{})
	config.DB.AutoMigrate(&entities.AccessToken{})
	config.DB.AutoMigrate(&entities.Follow{})

	userRepo := repositories.NewUserRepository(config.DB)
	followRepo := repositories.NewFollowRepository(config.DB)

	userUC := usecases.NewUserUsecase(userRepo)
	followUC := usecases.NewFollowUsecase(followRepo, userRepo)

	r := gin.Default()

	// Public routes
	r.POST("/register", httpUser.RegisterHandler(userUC))
	r.POST("/login", httpUser.LoginHandler(userUC))

	// Protected routes
	auth := r.Group("/auth")
	auth.Use(middleware.AuthMiddleware())
	{
		auth.GET("/profile", httpUser.ProfileHandler(userUC, followUC))
		auth.POST("/follow/:user_id", httpUser.FollowHandler(followUC))
		auth.DELETE("/follow/:user_id", httpUser.UnfollowHandler(followUC))
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
