package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/jstnangrendo/instagram-clone/like-service/config"
	"github.com/jstnangrendo/instagram-clone/like-service/domains/likes/entities"
	"github.com/jstnangrendo/instagram-clone/like-service/domains/likes/handlers/http"
	"github.com/jstnangrendo/instagram-clone/like-service/domains/likes/repositories"
	"github.com/jstnangrendo/instagram-clone/like-service/domains/likes/usecases"
	"github.com/jstnangrendo/instagram-clone/like-service/infrastructure/rabbitmq"
	"github.com/jstnangrendo/instagram-clone/like-service/middlewares"
)

func main() {
	// PostgreSQL
	config.InitDatabase()
	config.DB.AutoMigrate(&entities.Like{})

	// Redis
	config.InitRedis()

	// RabbitMQ
	config.InitRabbitMQ()
	defer config.RabbitConn.Close()

	likeRepo := repositories.NewLikeRepository(config.DB, config.RedisClient)
	publisher, err := rabbitmq.NewPublisherFromConn(config.RabbitConn, "post_like_queue")
	if err != nil {
		log.Fatalf("unable to create like-event publisher: %v", err)
	}
	defer publisher.Close()

	uc := usecases.NewLikeUseCase(likeRepo, publisher)
	h := http.NewLikeHandler(uc)

	r := gin.Default()
	r.Use(middlewares.AuthMiddleware())

	r.POST("/posts/:post_id/like", h.LikePost)
	r.DELETE("/posts/:post_id/unlike", h.UnlikePost)
	r.GET("/posts/:post_id/likes", h.GetLikeCount)
	r.GET("/popular", h.GetPopularPosts)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8085"
	}
	log.Println("Like service is running on port:" + port)
	r.Run(":" + port)
}
