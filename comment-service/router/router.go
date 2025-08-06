package router

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	commentHttp "github.com/jstnangrendo/instagram-clone/post-service/domains/comments/handlers/http"
	"github.com/jstnangrendo/instagram-clone/post-service/domains/comments/repositories"
	"github.com/jstnangrendo/instagram-clone/post-service/domains/comments/usecases"
	"github.com/jstnangrendo/instagram-clone/post-service/middlewares"
)

func NewRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	commentRepo := repositories.NewCommentRepository(db)
	commentUC := usecases.NewCommentUseCase(commentRepo)

	authGroup := r.Group("/")
	authGroup.Use(middlewares.AuthMiddleware())

	authGroup.POST("/posts/:post_id/comments", commentHttp.CreateCommentHandler(commentUC))
	r.GET("/posts/:post_id/comments", commentHttp.GetCommentsHandler(commentUC))

	return r
}
