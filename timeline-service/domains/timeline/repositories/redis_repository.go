package repositories

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/jstnangrendo/instagram-clone/timeline-service/domains/timeline/entities"
)

type Repository interface {
	AddToTimeline(ctx context.Context, userID int64, entry entities.PostCreatedEvent) error
	GetTimeline(ctx context.Context, userID int64, limit int64) ([]string, error)
}

type RedisRepository struct {
	client *redis.Client
}

func NewRedisRepository(client *redis.Client) *RedisRepository {
	return &RedisRepository{client: client}
}

func (r *RedisRepository) AddToTimeline(ctx context.Context, userID int64, entry entities.PostCreatedEvent) error {
	key := fmt.Sprintf("timeline:%d", userID)
	return r.client.ZAdd(ctx, key, &redis.Z{
		Score:  float64(entry.Timestamp.Unix()),
		Member: entry.PostID,
	}).Err()
}

func (r *RedisRepository) GetTimeline(ctx context.Context, userID int64, limit int64) ([]string, error) {
	key := fmt.Sprintf("timeline:%d", userID)
	return r.client.ZRevRange(ctx, key, 0, limit-1).Result()
}
