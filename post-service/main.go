package main

import (
	"log"

	"github.com/jstnangrendo/instagram-clone/post-service/config"
	"github.com/jstnangrendo/instagram-clone/post-service/domains/posts/entities"
	"github.com/jstnangrendo/instagram-clone/post-service/router"
)

func main() {
	config.InitDatabase()
	config.DB.AutoMigrate(&entities.Post{}, &entities.PostLike{})

	r := router.NewRouter(config.DB)

	log.Println("Post service running on :8081")
	r.Run(":8081")
}
