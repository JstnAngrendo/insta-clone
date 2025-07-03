package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jstnangrendo/instagram-clone/user-service/config"
	"github.com/jstnangrendo/instagram-clone/user-service/domains/users/entities"
	httpHandler "github.com/jstnangrendo/instagram-clone/user-service/domains/users/handlers/http"
	"github.com/jstnangrendo/instagram-clone/user-service/domains/users/repositories"
	"github.com/jstnangrendo/instagram-clone/user-service/domains/users/usecases"
	"github.com/jstnangrendo/instagram-clone/user-service/middleware"
)

func main() {
	// Init database & Redis
	config.InitDatabase()
	config.InitRedis()

	// Auto‑migrate tables
	config.DB.AutoMigrate(&entities.User{})
	config.DB.AutoMigrate(&entities.AccessToken{})

	// Setup repository & usecase
	userRepo := repositories.NewUserRepository(config.DB)
	userUC := usecases.NewUserUsecase(userRepo)

	// Setup Gin router
	r := gin.Default()

	// Public routes
	r.POST("/register", httpHandler.RegisterHandler(userUC))
	r.POST("/login", httpHandler.LoginHandler(userUC))

	// Protected routes group
	auth := r.Group("/auth")
	auth.Use(middleware.AuthMiddleware())
	auth.GET("/profile", httpHandler.ProfileHandler(userUC))

	// Run server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
