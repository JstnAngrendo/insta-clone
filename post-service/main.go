package main

import (
	"log"
	"os"

	"github.com/jstnangrendo/instagram-clone/post-service/config"
	"github.com/jstnangrendo/instagram-clone/post-service/domains/posts/entities"
	"github.com/jstnangrendo/instagram-clone/post-service/router"
)

func main() {

	err := os.MkdirAll("uploads", os.ModePerm)
	if err != nil {
		log.Fatalf("failed to create uploads directory: %v", err)
	}
	config.InitDatabase()
	config.DB.AutoMigrate(&entities.Post{}, &entities.PostLike{}, &entities.Tag{})

	r := router.NewRouter(config.DB)

	log.Println("Post service running on :8081")
	r.Run(":8081")
}
