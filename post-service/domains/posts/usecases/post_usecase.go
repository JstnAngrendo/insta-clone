package usecases

import (
	"fmt"
	"log"
	"time"

	"github.com/jstnangrendo/instagram-clone/post-service/domains/posts/entities"
	"github.com/jstnangrendo/instagram-clone/post-service/domains/posts/repositories"
	"github.com/jstnangrendo/instagram-clone/post-service/infrastructure/rabbitmq"
)

type PostUsecase interface {
	Create(userID uint, caption, imageURL string) (*entities.Post, error)
	GetByID(postID uint) (*entities.Post, error)
	GetByUser(userID uint) ([]entities.Post, error)
	Delete(postID, userID uint) error

	CreateWithTags(userID uint, caption, imageURL, thumbnailURL string, tags []entities.Tag) (*entities.Post, error)
	GetPostsByTag(tagName string) ([]entities.Post, error)

	GetByUserPaginated(userID uint, page int, size int) ([]entities.Post, int64, error)
	GetPostsByTagPaginated(tagName string, page int, size int) ([]entities.Post, int64, error)
	GetTimeline(userID string) ([]entities.Post, error)
}

type postUC struct {
	repo      repositories.PostRepository
	publisher *rabbitmq.Publisher
}

func NewPostUseCase(repo repositories.PostRepository, publisher *rabbitmq.Publisher) PostUsecase {
	return &postUC{repo: repo, publisher: publisher}
}

func (u *postUC) Create(userID uint, caption, imageURL string) (*entities.Post, error) {
	post := &entities.Post{
		UserID:    userID,
		Caption:   caption,
		ImageURL:  imageURL,
		CreatedAt: time.Now(),
	}
	if err := u.repo.Create(post); err != nil {
		return nil, err
	}
	event := map[string]interface{}{
		"post_id": post.ID,
		"user_id": post.UserID,
		"caption": post.Caption,
	}
	if err := u.publisher.Publish(event); err != nil {
		log.Printf("Failed to publish post_created event: %v", err)
	} else {
		fmt.Println("Published post_created event:", event)
	}
	return post, nil
}

func (u *postUC) GetByID(postID uint) (*entities.Post, error) {
	post, err := u.repo.FindByID(postID)
	return post, err
}

func (u *postUC) GetByUser(userID uint) ([]entities.Post, error) {
	return u.repo.FindByUserID(userID)
}

func (u *postUC) Delete(postID, userID uint) error {
	return u.repo.Delete(postID, userID)
}

func (u *postUC) CreateWithTags(userID uint, caption, imageURL, thumbnailURL string, tags []entities.Tag) (*entities.Post, error) {
	post := &entities.Post{
		UserID:       userID,
		Caption:      caption,
		ImageURL:     imageURL,
		ThumbnailURL: thumbnailURL,
		CreatedAt:    time.Now(),
	}
	if err := u.repo.CreateWithTags(post, tags); err != nil {
		return nil, err
	}
	event := map[string]interface{}{
		"post_id": post.ID,
		"user_id": post.UserID,
		"caption": post.Caption,
	}
	if err := u.publisher.Publish(event); err != nil {
		log.Printf("Failed to publish post_created event: %v", err)
	} else {
		fmt.Println("Published post_created event:", event)
	}
	return post, nil
}

func (u *postUC) GetPostsByTag(tagName string) ([]entities.Post, error) {
	return u.repo.FindPostsByTag(tagName)
}

func (u *postUC) GetByUserPaginated(userID uint, page int, size int) ([]entities.Post, int64, error) {
	offset := (page - 1) * size
	var posts []entities.Post
	var total int64

	err := u.repo.CountByUser(userID, &total)
	if err != nil {
		return nil, 0, err
	}

	err = u.repo.FindByUserPaginated(userID, offset, size, &posts)
	if err != nil {
		return nil, 0, err
	}
	return posts, total, nil
}

func (u *postUC) GetPostsByTagPaginated(tagName string, page int, size int) ([]entities.Post, int64, error) {
	offset := (page - 1) * size
	var posts []entities.Post
	var total int64

	err := u.repo.CountPostsByTag(tagName, &total)
	if err != nil {
		return nil, 0, err
	}

	err = u.repo.FindPostsByTagPaginated(tagName, offset, size, &posts)
	if err != nil {
		return nil, 0, err
	}
	return posts, total, nil
}

func (u *postUC) GetTimeline(userID string) ([]entities.Post, error) {
	return u.repo.GetPostsByUserIDs([]string{userID})
}
