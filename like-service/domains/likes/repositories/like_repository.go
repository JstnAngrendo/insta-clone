package repositories

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/jstnangrendo/instagram-clone/like-service/domains/likes/entities"
	"gorm.io/gorm"
)

type LikeRepository struct {
	db  *gorm.DB
	rdb *redis.Client
}

func NewLikeRepository(db *gorm.DB, rdb *redis.Client) *LikeRepository {
	return &LikeRepository{db: db, rdb: rdb}
}

func (r *LikeRepository) LikePost(ctx context.Context, postID, userID uint) error {
	// check if already liked
	var count int64
	err := r.db.Model(&entities.Like{}).
		Where("post_id = ? AND user_id = ?", postID, userID).
		Count(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		// already liked, no-op
		return nil
	}


	like := &entities.Like{PostID: postID, UserID: userID}
	if err := r.db.Create(like).Error; err != nil {
		return err
	}

	// increment Redis sorted set
	key := "popular_posts"
	if _, err := r.rdb.ZIncrBy(ctx, key, 1, fmt.Sprintf("%d", postID)).Result(); err != nil {
		return err
	}
	return nil
}


func (r *LikeRepository) UnlikePost(ctx context.Context, postID, userID uint) error {
	// delete from DB
	tx := r.db.Where("post_id = ? AND user_id = ?", postID, userID).Delete(&entities.Like{})
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		// nothing to unlike
		return nil
	}

	// decrement Redis sorted set
	key := "popular_posts"
	if _, err := r.rdb.ZIncrBy(ctx, key, -1, fmt.Sprintf("%d", postID)).Result(); err != nil {
		return err
	}
	return nil
}


func (r *LikeRepository) GetLikeCount(ctx context.Context, postID uint) (int64, error) {
	var count int64
	err := r.db.Model(&entities.Like{}).
		Where("post_id = ?", postID).
		Count(&count).Error
	return count, err
}

// GetPopularPosts returns top N popular post IDs from Redis
func (r *LikeRepository) GetPopularPosts(ctx context.Context, limit int64) ([]uint, error) {
	key := "popular_posts"
	members, err := r.rdb.ZRevRange(ctx, key, 0, limit-1).Result()
	if err != nil {
		return nil, err
	}
	var ids []uint
	for _, m := range members {
		var id uint
		if _, err := fmt.Sscan(m, &id); err == nil {
			ids = append(ids, id)
		}
	}
	return ids, nil
}
