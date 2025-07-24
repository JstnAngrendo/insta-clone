package repositories

import (
	"github.com/jstnangrendo/instagram-clone/post-service/domains/comments/entities"
	"gorm.io/gorm"
)

type commentRepository struct {
	db *gorm.DB
}

type CommentRepository interface {
	SaveComment(comment *entities.Comment) error
	FetchComments(postID uint) ([]entities.Comment, error)
}

func NewCommentRepository(db *gorm.DB) CommentRepository {
	return &commentRepository{db}
}

func (r *commentRepository) SaveComment(comment *entities.Comment) error {
	return r.db.Create(comment).Error
}

func (r *commentRepository) FetchComments(postID uint) ([]entities.Comment, error) {
	var comments []entities.Comment
	err := r.db.Where("post_id = ?", postID).Order("created_at asc").Find(&comments).Error
	return comments, err
}
