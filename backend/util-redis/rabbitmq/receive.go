package rabbitmq

import (
	"context"
	"strings"
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
			data := strings.Split(string(d.Body), "|")
			if len(data) != 2 {
				zap.L().Error("Invalid message format.")
				continue
			}

			zap.L().Info("Setting Redis key: " + data[1])

			err = redis_client.Client.SetEx(
				context.Background(),
				data[0],
				data[1],
				time.Duration(60)*time.Second,
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
