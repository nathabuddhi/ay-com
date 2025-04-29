package redis_client

import (
	"context"
	"os"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var Client *redis.Client

func InitRedis() {
	addr := os.Getenv("REDIS_URL")

	Client = redis.NewClient(&redis.Options{
		Addr: addr,
		DB:   0,
	})

	if err := Client.Ping(context.Background()).Err(); err != nil {
		zap.L().Error("Failed to connect to Redis: " + err.Error())
	}
	zap.L().Info("Redis client initialized successfully.")
}

func GetCache(key string) string {
	return Client.Get(context.Background(), key).Val()
}
