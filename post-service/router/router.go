package router

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	postHttp "github.com/jstnangrendo/instagram-clone/post-service/domains/posts/handlers/http"
	"github.com/jstnangrendo/instagram-clone/post-service/domains/posts/repositories"
	"github.com/jstnangrendo/instagram-clone/post-service/domains/posts/usecases"
	"github.com/jstnangrendo/instagram-clone/post-service/middlewares"
)

func NewRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	postRepo := repositories.NewPostRepository(db)
	postUC := usecases.NewPostUsecase(postRepo)

	authGroup := r.Group("/")
	authGroup.Use(middlewares.AuthMiddleware())

	authGroup.POST("/posts", postHttp.CreatePostHandler(postUC))
	authGroup.DELETE("/posts/:post_id", postHttp.DeletePostHandler(postUC))

	r.GET("/posts/:post_id", postHttp.GetPostHandler(postUC))
	authGroup.GET("/users/:user_id/posts", postHttp.GetUserPostsHandler(postUC))

	authGroup.POST("/posts/:post_id/like", postHttp.LikePostHandler(postUC))
	authGroup.DELETE("/posts/:post_id/unlike", postHttp.UnlikePostHandler(postUC))

	return r
}
