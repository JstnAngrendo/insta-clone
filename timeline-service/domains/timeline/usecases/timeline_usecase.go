package usecases

import (
	"context"

	"github.com/jstnangrendo/instagram-clone/timeline-service/domains/timeline/entities"
	"github.com/jstnangrendo/instagram-clone/timeline-service/domains/timeline/repositories"
)

type TimelineUseCase interface {
	ProcessNewPost(ctx context.Context, evt entities.PostCreatedEvent) error
	FetchTimeline(ctx context.Context, userID int64, limit int64) ([]string, error)
}

type timelineUseCase struct {
	repo            repositories.Repository
	followerService FollowerService
}

type FollowerService interface {
	GetFollowers(ctx context.Context, userID int64) ([]int64, error)
}

func NewTimelineUseCase(repo repositories.Repository, followerService FollowerService) TimelineUseCase {
	return &timelineUseCase{
		repo:            repo,
		followerService: followerService,
	}
}

func (u *timelineUseCase) ProcessNewPost(ctx context.Context, evt entities.PostCreatedEvent) error {
	followers, err := u.followerService.GetFollowers(ctx, evt.AuthorID)
	if err != nil {
		return err
	}

	for _, followerID := range followers {
		err := u.repo.AddToTimeline(ctx, followerID, evt)
		if err != nil {
			continue
		}
	}

	return nil
}

func (u *timelineUseCase) FetchTimeline(ctx context.Context, userID int64, limit int64) ([]string, error) {
	return u.repo.GetTimeline(ctx, userID, limit)
}
