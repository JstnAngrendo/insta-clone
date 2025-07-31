package main

import (
	"log"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/jstnangrendo/instagram-clone/timeline-service/client"
	"github.com/jstnangrendo/instagram-clone/timeline-service/config"
	"github.com/jstnangrendo/instagram-clone/timeline-service/domains/timeline/entities"
	"github.com/jstnangrendo/instagram-clone/timeline-service/domains/timeline/repositories"
	"github.com/jstnangrendo/instagram-clone/timeline-service/domains/timeline/usecases"
	"github.com/jstnangrendo/instagram-clone/timeline-service/handlers"
	"github.com/jstnangrendo/instagram-clone/timeline-service/infrastructure/rabbitmq"
	"github.com/jstnangrendo/instagram-clone/timeline-service/router"
)

func main() {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	err := os.MkdirAll("uploads", os.ModePerm)
	if err != nil {
		log.Fatalf("failed to create uploads directory: %v", err)
	}

	config.InitDatabase()
	config.DB.AutoMigrate(&entities.Timeline{})
	timelineRepo := repositories.NewRedisRepository(redisClient)
	followerService := client.NewFollowerService("http://localhost:8080")

	timelineUsecase := usecases.NewTimelineUseCase(timelineRepo, followerService)

	amqpURL := "amqp://guest:guest@localhost:5672/"
	consumer, err := rabbitmq.NewConsumer(amqpURL, timelineUsecase)
	if err != nil {
		log.Fatalf("Failed to create RabbitMQ consumer: %v", err)
	}
	defer consumer.Close()

	err = consumer.StartConsuming("post_created_queue")
	if err != nil {
		log.Fatalf("Failed to start consuming messages: %v", err)
	}

	postCli := client.NewPostClient("http://localhost:8081")
	handler := handlers.NewTimelineHandler(timelineUsecase, postCli)
	httpSrv := router.SetupRouter(handler)
	go func() {
		log.Println("HTTP server running on :8084")
		if err := httpSrv.Run(":8084"); err != nil {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	select {}
}
