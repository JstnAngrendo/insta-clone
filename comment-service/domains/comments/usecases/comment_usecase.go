package usecases

import (
	"github.com/jstnangrendo/instagram-clone/post-service/domains/comments/entities"
	"github.com/jstnangrendo/instagram-clone/post-service/domains/comments/repositories"
)

type CommentUsecase interface {
	CreateComment(postID uint, userID uint, content string) error
	GetComments(postID uint) ([]entities.Comment, error)
}

type commentUC struct {
	repo repositories.CommentRepository
}

func NewCommentUseCase(commentRepo repositories.CommentRepository) CommentUsecase {
	return &commentUC{
		repo: commentRepo,
	}
}

func (u *commentUC) CreateComment(postID uint, userID uint, content string) error {
	comment := entities.Comment{
		PostID:  postID,
		UserID:  userID,
		Content: content,
	}
	return u.repo.SaveComment(&comment)
}

func (u *commentUC) GetComments(postID uint) ([]entities.Comment, error) {
	return u.repo.FetchComments(postID)
}
