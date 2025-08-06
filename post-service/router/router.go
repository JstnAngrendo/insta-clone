package router

import (
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	postHttp "github.com/jstnangrendo/instagram-clone/post-service/domains/posts/handlers/http"
	"github.com/jstnangrendo/instagram-clone/post-service/domains/posts/repositories"
	usecases "github.com/jstnangrendo/instagram-clone/post-service/domains/posts/usecases"
	"github.com/jstnangrendo/instagram-clone/post-service/infrastructure/rabbitmq"
	"github.com/jstnangrendo/instagram-clone/post-service/middlewares"
)

func NewRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	postRepo := repositories.NewPostRepository(db)

	publisher, err := rabbitmq.NewPublisher("amqp://guest:guest@localhost:5672/", "post_created_queue")
	if err != nil {
		log.Fatalf("Failed to create RabbitMQ publisher: %v", err)
	}

	postUC := usecases.NewPostUseCase(postRepo, publisher)

	h := postHttp.NewPostHandler(postUC)

	r.GET("/posts/:post_id", h.GetByID)
	r.GET("/posts/tag/:tagName", h.GetByTag)
	r.POST("/posts/batch", postHttp.BatchGetHandler(postUC))

	auth := r.Group("/")
	auth.Use(middlewares.AuthMiddleware())
	{
		auth.POST("/posts", h.Create)
		auth.DELETE("/posts/:post_id", h.Delete)
		auth.GET("/users/:user_id/posts", h.GetByUser)
		auth.GET("/timeline", h.Timeline)
	}

	return r
}
