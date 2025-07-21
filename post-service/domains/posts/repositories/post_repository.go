package repositories

import (
	"github.com/jstnangrendo/instagram-clone/post-service/domains/posts/entities"
	"gorm.io/gorm"
)

type postRepository struct {
	db *gorm.DB
}

type PostRepository interface {
	Create(post *entities.Post) error
	FindByID(id uint) (*entities.Post, error)
	FindByUserID(userID uint) ([]entities.Post, error)
	Delete(id uint, userID uint) error

	LikePost(userID, postID uint) error
	UnlikePost(userID, postID uint) error
	CountLikes(postID uint) (int64, error)
}

func NewPostRepository(db *gorm.DB) PostRepository {
	return &postRepository{db}
}

func (r *postRepository) Create(post *entities.Post) error {
	return r.db.Create(post).Error
}

func (r *postRepository) FindByID(id uint) (*entities.Post, error) {
	var post entities.Post
	err := r.db.First(&post, id).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *postRepository) FindByUserID(userID uint) ([]entities.Post, error) {
	var posts []entities.Post
	err := r.db.Where("user_id = ?", userID).Order("created_at desc").Find(&posts).Error
	return posts, err
}

func (r *postRepository) Delete(id uint, userID uint) error {
	return r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&entities.Post{}).Error
}

func (r *postRepository) LikePost(userID, postID uint) error {
	like := entities.PostLike{
		UserID: userID,
		PostID: postID,
	}
	return r.db.Create(&like).Error
}

func (r *postRepository) UnlikePost(userID, postID uint) error {
	return r.db.Where("user_id = ? AND post_id = ?", userID, postID).Delete(&entities.PostLike{}).Error
}

func (r *postRepository) CountLikes(postID uint) (int64, error) {
	var count int64
	err := r.db.Model(&entities.PostLike{}).
		Where("post_id = ?", postID).
		Count(&count).Error
	return count, err
}
