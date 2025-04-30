package redis_client

import (
	"context"
	"encoding/json"
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

func GetCache(key string) interface{} {
	value, err := Client.Get(context.Background(), key).Result()
	if err != nil || value == "" {
		return nil
	}

	var payload interface{}
	if value != "" {
		if err := json.Unmarshal([]byte(value), &payload); err != nil {
			zap.L().Error("Failed to unmarshal cached JSON", zap.String("key", value), zap.Error(err))
			return nil
		}
		zap.L().Info("User profile fetched from Redis", zap.Any("payload", payload))
		return payload
	}
	return nil
}
