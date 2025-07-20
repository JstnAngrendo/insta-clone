package repositories

import (
	"errors"

	"github.com/jstnangrendo/instagram-clone/user-service/domains/users/entities"
	"gorm.io/gorm"
)

type FollowRepository interface {
	Follow(followerID, followeeID uint) error
	Unfollow(followerID, followeeID uint) error
	IsFollowing(followerID, followeeID uint) (bool, error)
	CountFollowers(userID uint) (int64, error)
	CountFollowing(userID uint) (int64, error)
}

type followRepo struct {
	db *gorm.DB
}

func NewFollowRepository(db *gorm.DB) FollowRepository {
	return &followRepo{db}
}

func (r *followRepo) Follow(followerID, followeeID uint) error {
	return r.db.Create(&entities.Follow{FollowerID: followerID, FolloweeID: followeeID}).Error
}

func (r *followRepo) Unfollow(followerID, followeeID uint) error {
	return r.db.
		Where("follower_id = ? AND followee_id = ?", followerID, followeeID).
		Delete(&entities.Follow{}).Error
}

func (r *followRepo) IsFollowing(followerID, followeeID uint) (bool, error) {
	var f entities.Follow
	err := r.db.
		Where("follower_id = ? AND followee_id = ?", followerID, followeeID).
		First(&f).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (r *followRepo) CountFollowers(userID uint) (int64, error) {
	var c int64
	err := r.db.Model(&entities.Follow{}).
		Where("followee_id = ?", userID).
		Count(&c).Error
	return c, err
}

func (r *followRepo) CountFollowing(userID uint) (int64, error) {
	var c int64
	err := r.db.Model(&entities.Follow{}).
		Where("follower_id = ?", userID).
		Count(&c).Error
	return c, err
}
