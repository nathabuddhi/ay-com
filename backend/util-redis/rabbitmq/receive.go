package rabbitmq

import (
	"context"
	"encoding/json"
	"time"

	"github.com/nathabuddhi/ay-com/backend/util-redis/redis_client"
	"github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

var SetRedisChannel *amqp091.Channel
var DeleteRedisChannel *amqp091.Channel

var Forever = make(chan bool)

func InitSetRedisChannel(conn *amqp091.Connection) {
	ch, err := conn.Channel()
	if err != nil {
		zap.L().Fatal("Failed to open a channel: " + err.Error())
	}

	_, err = ch.QueueDeclare(
		"set_redis",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		zap.L().Fatal("Failed to declare set_redis queue: " + err.Error())
	} else {
		zap.L().Info("RabbitMQ set_redis queue declared successfully")
	}

	SetRedisChannel = ch
}

func InitDeleteRedisChannel(conn *amqp091.Connection) {
	ch, err := conn.Channel()
	if err != nil {
		zap.L().Fatal("Failed to open a channel: " + err.Error())
	}

	_, err = ch.QueueDeclare(
		"delete_redis",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		zap.L().Fatal("Failed to declare delete_redis queue: " + err.Error())
	} else {
		zap.L().Info("RabbitMQ delete_redis queue declared successfully")
	}

	DeleteRedisChannel = ch
}

type RedisPayload struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	TTL   int    `json:"ttl"`
}

func StartConsumingSet() {
	msgs, err := SetRedisChannel.Consume(
		"set_redis",
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		zap.L().Fatal("Failed to register consumer on set_redis: " + err.Error())
	}

	go func() {
		for d := range msgs {
			var payload RedisPayload
			err := json.Unmarshal(d.Body, &payload)
			if err != nil {
				zap.L().Error("Failed to parse set_redis message: " + err.Error())
				continue
			}

			zap.L().Info("Setting Redis key: " + payload.Key)

			err = redis_client.Client.SetEx(
				context.Background(),
				payload.Key,
				payload.Value,
				time.Duration(payload.TTL)*time.Second,
			).Err()

			if err != nil {
				zap.L().Error("Redis SET error: " + err.Error())
			}
		}
	}()
}

func StartConsumingDelete() {
	msgs, err := DeleteRedisChannel.Consume(
		"delete_redis",
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		zap.L().Fatal("Failed to register consumer on delete_redis: " + err.Error())
	}

	go func() {
		for d := range msgs {
			key := string(d.Body)
			zap.L().Info("Deleting Redis key: " + key)

			err := redis_client.Client.Del(context.Background(), key).Err()
			if err != nil {
				zap.L().Error("Redis DEL error: " + err.Error())
			}
		}
	}()

	zap.L().Info(" [*] Waiting for messages...")
}
