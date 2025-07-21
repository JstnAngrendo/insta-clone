package usecases

import (
	"errors"

	"github.com/jstnangrendo/instagram-clone/user-service/domains/users/repositories"
)

type FollowUsecase interface {
	Follow(userID, targetID uint) error
	Unfollow(userID, targetID uint) error
	GetCounts(userID uint) (followers, following int64, err error)
	GetFollowingUserIDs(userID uint) ([]uint, error)
}

type followUC struct {
	repo     repositories.FollowRepository
	userRepo repositories.UserRepository
}

func NewFollowUsecase(
	fRepo repositories.FollowRepository,
	uRepo repositories.UserRepository,
) FollowUsecase {
	return &followUC{
		repo:     fRepo,
		userRepo: uRepo,
	}
}

func (u *followUC) Follow(userID, targetID uint) error {
	if userID == targetID {
		return errors.New("cannot follow yourself")
	}

	if _, err := u.userRepo.FindByID(targetID); err != nil {
		return errors.New("target user not found")
	}

	already, err := u.repo.IsFollowing(userID, targetID)
	if err != nil {
		return err
	}
	if already {
		return errors.New("already following this user")
	}

	return u.repo.Follow(userID, targetID)
}

func (u *followUC) Unfollow(userID, targetID uint) error {
	already, err := u.repo.IsFollowing(userID, targetID)
	if err != nil {
		return err
	}
	if !already {
		return errors.New("you are not following this user")
	}
	return u.repo.Unfollow(userID, targetID)
}

func (u *followUC) GetCounts(userID uint) (followers, following int64, err error) {
	followers, err = u.repo.CountFollowers(userID)
	if err != nil {
		return
	}
	following, err = u.repo.CountFollowing(userID)
	return
}

func (u *followUC) GetFollowingUserIDs(userID uint) ([]uint, error) {
	return u.repo.GetFollowingUserIDs(userID)
}
