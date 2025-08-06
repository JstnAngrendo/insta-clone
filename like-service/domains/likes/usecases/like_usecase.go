package usecases

import (
	"context"
	"fmt"

	"github.com/jstnangrendo/instagram-clone/like-service/domains/likes/repositories"
	"github.com/jstnangrendo/instagram-clone/like-service/infrastructure/rabbitmq"
)

type LikeUseCase struct {
	repo      *repositories.LikeRepository
	publisher *rabbitmq.Publisher
}

type PostLikeEvent struct {
	PostID uint   `json:"post_id"`
	Action string `json:"action"`
}

func NewLikeUseCase(repo *repositories.LikeRepository, pub *rabbitmq.Publisher) *LikeUseCase {
	return &LikeUseCase{
		repo:      repo,
		publisher: pub,
	}
}

func (uc *LikeUseCase) LikePost(ctx context.Context, postID, userID uint) error {
	if err := uc.repo.LikePost(ctx, postID, userID); err != nil {
		return err
	}
	evt := PostLikeEvent{PostID: postID, Action: "like"}
	if err := uc.publisher.Publish(evt); err != nil {
		fmt.Printf("[LikeUseCase] failed to publish like event: %v\n", err)
	}
	return nil
}

func (uc *LikeUseCase) UnlikePost(ctx context.Context, postID, userID uint) error {
	if err := uc.repo.UnlikePost(ctx, postID, userID); err != nil {
		return err
	}
	evt := PostLikeEvent{PostID: postID, Action: "unlike"}
	if err := uc.publisher.Publish(evt); err != nil {
		fmt.Printf("[LikeUseCase] failed to publish unlike event: %v\n", err)
	}
	return nil
}

func (uc *LikeUseCase) GetLikeCount(ctx context.Context, postID uint) (int64, error) {
	return uc.repo.GetLikeCount(ctx, postID)
}

func (uc *LikeUseCase) GetPopularPosts(ctx context.Context, limit int64) ([]uint, error) {
	return uc.repo.GetPopularPosts(ctx, limit)
}
