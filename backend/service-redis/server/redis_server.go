package server

import (
	"context"
	"time"

	"github.com/nathabuddhi/ay-com/backend/service-redis/redis_client"
	"github.com/redis/go-redis/v9"

	pb "github.com/nathabuddhi/ay-com/backend/service-redis/proto"
)

type RedisServer struct {
	pb.UnimplementedRedisServiceServer
}

func NewRedisServer() *RedisServer {
	return &RedisServer{}
}

func (s *RedisServer) SetKey(ctx context.Context, req *pb.SetKeyRequest) (*pb.SetKeyResponse, error) {
	err := redis_client.Client.Set(ctx, req.Key, req.Value, time.Duration(req.ExpirationSeconds)*time.Second).Err()
	if err != nil {
		return &pb.SetKeyResponse{Success: false}, err
	}
	return &pb.SetKeyResponse{Success: true}, nil
}

func (s *RedisServer) GetKey(ctx context.Context, req *pb.GetKeyRequest) (*pb.GetKeyResponse, error) {
	val, err := redis_client.Client.Get(ctx, req.Key).Result()
	if err == redis.Nil {
		return &pb.GetKeyResponse{Found: false}, nil
	} else if err != nil {
		return nil, err
	}
	return &pb.GetKeyResponse{Value: val, Found: true}, nil
}

func (s *RedisServer) DeleteKey(ctx context.Context, req *pb.DeleteKeyRequest) (*pb.DeleteKeyResponse, error) {
	err := redis_client.Client.Del(ctx, req.Key).Err()
	if err != nil {
		return &pb.DeleteKeyResponse{Success: false}, err
	}
	return &pb.DeleteKeyResponse{Success: true}, nil
}
