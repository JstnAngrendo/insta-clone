package main

import (
	"log"

	"github.com/jstnangrendo/instagram-clone/post-service/config"
	"github.com/jstnangrendo/instagram-clone/post-service/domains/comments/entities"
	"github.com/jstnangrendo/instagram-clone/post-service/router"
)

func main() {

	config.InitDatabase()
	config.DB.AutoMigrate(&entities.Comment{})

	r := router.NewRouter(config.DB)

	log.Println("Post service running on :8082")
	r.Run(":8082")
}
