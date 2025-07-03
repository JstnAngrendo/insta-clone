package config

import (
	"context"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client
var Ctx = context.Background()

func InitRedis() {
	addr := os.Getenv("REDIS_ADDR")
	RedisClient = redis.NewClient(&redis.Options{Addr: addr})
	if err := RedisClient.Ping(Ctx).Err(); err != nil {
		log.Fatal("Redis connection failed:", err)
	}
}
