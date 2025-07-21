package usecases

import (
	"time"

	"github.com/jstnangrendo/instagram-clone/post-service/domains/posts/entities"
	"github.com/jstnangrendo/instagram-clone/post-service/domains/posts/repositories"
)

type PostUsecase interface {
	Create(userID uint, caption, imageURL string) (*entities.Post, error)
	GetByID(postID uint) (*entities.Post, error)
	GetByUser(userID uint) ([]entities.Post, error)
	Delete(postID, userID uint) error
	LikePost(userID, postID uint) error
	UnlikePost(userID, postID uint) error
	CountLikes(postID uint) (int64, error)
	CreateWithTags(userID uint, caption, imageURL, thumbnailURL string, tags []entities.Tag) (*entities.Post, error)
	GetPostsByTag(tagName string) ([]entities.Post, error)

	GetByUserPaginated(userID uint, page int, size int) ([]entities.Post, int64, error)
	GetPostsByTagPaginated(tagName string, page int, size int) ([]entities.Post, int64, error)
}

type postUC struct {
	repo repositories.PostRepository
}

func NewPostUsecase(r repositories.PostRepository) PostUsecase {
	return &postUC{repo: r}
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
	return post, nil
}

func (u *postUC) GetByID(postID uint) (*entities.Post, error) {
	post, err := u.repo.FindByID(postID)
	if err != nil {
		return nil, err
	}

	likeCount, err := u.repo.CountLikes(postID)
	if err != nil {
		return nil, err
	}
	post.LikeCount = likeCount

	return post, nil
}

func (u *postUC) GetByUser(userID uint) ([]entities.Post, error) {
	posts, err := u.repo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	for i := range posts {
		likeCount, err := u.repo.CountLikes(posts[i].ID)
		if err != nil {
			return nil, err
		}
		posts[i].LikeCount = likeCount
	}

	return posts, nil
}

func (u *postUC) Delete(postID, userID uint) error {
	return u.repo.Delete(postID, userID)
}

func (u *postUC) LikePost(userID, postID uint) error {
	return u.repo.LikePost(userID, postID)
}

func (u *postUC) UnlikePost(userID, postID uint) error {
	return u.repo.UnlikePost(userID, postID)
}

func (u *postUC) CountLikes(postID uint) (int64, error) {
	return u.repo.CountLikes(postID)
}

func (uc *postUC) CreateWithTags(userID uint, caption, imageURL, thumbnailURL string, tags []entities.Tag) (*entities.Post, error) {
	post := &entities.Post{
		UserID:       userID,
		Caption:      caption,
		ImageURL:     imageURL,
		ThumbnailURL: thumbnailURL,
		CreatedAt:    time.Now(),
	}

	err := uc.repo.CreateWithTags(post, tags)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (uc *postUC) GetPostsByTag(tagName string) ([]entities.Post, error) {
	return uc.repo.FindPostsByTag(tagName)
}

func (uc *postUC) GetByUserPaginated(userID uint, page int, size int) ([]entities.Post, int64, error) {
	offset := (page - 1) * size
	var posts []entities.Post
	var total int64

	err := uc.repo.CountByUser(userID, &total)
	if err != nil {
		return nil, 0, err
	}

	err = uc.repo.FindByUserPaginated(userID, offset, size, &posts)
	if err != nil {
		return nil, 0, err
	}

	for i := range posts {
		likeCount, err := uc.repo.CountLikes(posts[i].ID)
		if err != nil {
			return nil, 0, err
		}
		posts[i].LikeCount = likeCount
	}

	return posts, total, nil
}

func (uc *postUC) GetPostsByTagPaginated(tagName string, page int, size int) ([]entities.Post, int64, error) {
	offset := (page - 1) * size
	var posts []entities.Post
	var total int64

	err := uc.repo.CountPostsByTag(tagName, &total)
	if err != nil {
		return nil, 0, err
	}

	err = uc.repo.FindPostsByTagPaginated(tagName, offset, size, &posts)
	if err != nil {
		return nil, 0, err
	}

	for i := range posts {
		likeCount, err := uc.repo.CountLikes(posts[i].ID)
		if err != nil {
			return nil, 0, err
		}
		posts[i].LikeCount = likeCount
	}

	return posts, total, nil
}
