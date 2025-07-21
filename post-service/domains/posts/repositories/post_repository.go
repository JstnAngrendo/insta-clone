package repositories

import (
	"log"
	"strings"

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
	CreateWithTags(post *entities.Post, tags []entities.Tag) error
	FindPostsByTag(tagName string) ([]entities.Post, error)

	FindByUserPaginated(userID uint, offset int, limit int, posts *[]entities.Post) error
	CountByUser(userID uint, total *int64) error

	FindPostsByTagPaginated(tagName string, offset int, limit int, posts *[]entities.Post) error
	CountPostsByTag(tagName string, total *int64) error
}

func NewPostRepository(db *gorm.DB) PostRepository {
	return &postRepository{db}
}

func (r *postRepository) Create(post *entities.Post) error {
	return r.db.Create(post).Error
}

func (r *postRepository) FindByID(id uint) (*entities.Post, error) {
	var post entities.Post
	err := r.db.Preload("Tags").First(&post, id).Error
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

func (r *postRepository) CreateWithTags(post *entities.Post, tags []entities.Tag) error {
	for i := range tags {
		var existing entities.Tag
		if err := r.db.Where("name = ?", tags[i].Name).First(&existing).Error; err == nil {
			tags[i].ID = existing.ID
		} else {
			if err := r.db.Create(&tags[i]).Error; err != nil {
				return err
			}
		}
	}

	post.Tags = tags
	if err := r.db.Create(post).Error; err != nil {
		return err
	}

	return r.db.Preload("Tags").First(post, post.ID).Error
}

func (r *postRepository) FindPostsByTag(tagName string) ([]entities.Post, error) {
	var posts []entities.Post
	err := r.db.
		Joins("JOIN post_tags ON post_tags.post_id = posts.id").
		Joins("JOIN tags ON tags.id = post_tags.tag_id").
		Where("LOWER(tags.name) = ?", strings.ToLower(tagName)).
		Preload("Tags").
		Find(&posts).Error

	if err != nil {
		log.Println("DB error:", err)
	} else {
		log.Println("Number of posts found with tag", tagName, ":", len(posts))
	}
	return posts, err
}

func (r *postRepository) FindByUserPaginated(userID uint, offset int, limit int, posts *[]entities.Post) error {
	return r.db.Where("user_id = ?", userID).
		Order("created_at desc").
		Offset(offset).
		Limit(limit).
		Preload("Tags").
		Find(posts).Error
}

func (r *postRepository) CountByUser(userID uint, total *int64) error {
	return r.db.Model(&entities.Post{}).Where("user_id = ?", userID).Count(total).Error
}

func (r *postRepository) FindPostsByTagPaginated(tagName string, offset int, limit int, posts *[]entities.Post) error {
	return r.db.
		Joins("JOIN post_tags ON post_tags.post_id = posts.id").
		Joins("JOIN tags ON tags.id = post_tags.tag_id").
		Where("LOWER(tags.name) = ?", strings.ToLower(tagName)).
		Order("posts.created_at desc").
		Offset(offset).
		Limit(limit).
		Preload("Tags").
		Find(posts).Error
}

func (r *postRepository) CountPostsByTag(tagName string, total *int64) error {
	return r.db.
		Model(&entities.Post{}).
		Joins("JOIN post_tags ON post_tags.post_id = posts.id").
		Joins("JOIN tags ON tags.id = post_tags.tag_id").
		Where("LOWER(tags.name) = ?", strings.ToLower(tagName)).
		Count(total).Error
}
