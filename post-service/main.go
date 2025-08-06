package main

import (
	"log"
	"os"

	"github.com/jstnangrendo/instagram-clone/post-service/config"
	"github.com/jstnangrendo/instagram-clone/post-service/domains/posts/entities"
	"github.com/jstnangrendo/instagram-clone/post-service/domains/posts/repositories"
	"github.com/jstnangrendo/instagram-clone/post-service/infrastructure/rabbitmq"
	"github.com/jstnangrendo/instagram-clone/post-service/router"
)

func main() {
	er := os.MkdirAll("uploads", os.ModePerm)
	if er != nil {
		log.Fatalf("failed to create uploads directory: %v", er)
	}

	config.InitDatabase()
	config.DB.AutoMigrate(&entities.Post{}, &entities.Tag{})
	config.InitRabbitMQ()
	defer config.CloseRabbitMQ()

	postRepo := repositories.NewPostRepository(config.DB)
	likeConsumer, err := rabbitmq.NewConsumer(config.RabbitConn, postRepo)
	if err != nil {
		log.Fatalf("Failed to start like consumer: %v", err)
	}
	likeConsumer.Start("post_like_queue")

	r := router.NewRouter(config.DB)
	log.Println("Post service running on :8081")
	r.Run(":8081")
}
