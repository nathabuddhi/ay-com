package rabbitmq

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

var NotificationChannel *amqp.Channel
var SetRedisChannel *amqp.Channel
var DeleteRedisChannel *amqp.Channel
var SendMentionChannel *amqp.Channel
var SendRepostChannel *amqp.Channel

func InitMentionChannel(conn *amqp.Connection) {
	ch, err := conn.Channel()
	if err != nil {
		zap.L().Fatal("Failed to open a mention_user channel: " + err.Error())
	}

	_, err = ch.QueueDeclare(
		"mention_user",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		zap.L().Fatal("Failed to declare mention_user queue: " + err.Error())
	} else {
		zap.L().Info("RabbitMQ mention_user queue declared successfully")
	}

	SendMentionChannel = ch
}

func InitNotificationChannel(conn *amqp.Connection) {
	ch, err := conn.Channel()
	if err != nil {
		zap.L().Fatal("Failed to open a send_notification channel: " + err.Error())
	}

	_, err = ch.QueueDeclare(
		"send_notification",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		zap.L().Fatal("Failed to declare send_notification queue: " + err.Error())
	} else {
		zap.L().Info("RabbitMQ send_notification queue declared successfully")
	}

	NotificationChannel = ch
}

func InitSetRedisChannel(conn *amqp.Connection) {
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

func InitDeleteRedisChannel(conn *amqp.Connection) {
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

func InitRepostChannel(conn *amqp.Connection) {
	ch, err := conn.Channel()
	if err != nil {
		zap.L().Fatal("Failed to open a channel: " + err.Error())
	}

	_, err = ch.QueueDeclare(
		"send_repost",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		zap.L().Fatal("Failed to declare send_repost queue: " + err.Error())
	} else {
		zap.L().Info("RabbitMQ send_repost queue declared successfully")
	}

	SendRepostChannel = ch
}

func InitRabbitMQ() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		zap.L().Fatal("Failed to connect to RabbitMQ: " + err.Error())
	}

	InitSetRedisChannel(conn)
	InitDeleteRedisChannel(conn)
	InitNotificationChannel(conn)
	InitMentionChannel(conn)
	InitRepostChannel(conn)
}

func PublishSendNotification(notif_type string, user_id string, email string, title string, content string, from string) error {
	body := notif_type + "|" + user_id + "|" + email + "|" + title + "|" + content + "|" + from

	zap.L().Info("Publishing send_notification message to RabbitMQ: " + body)

	err := NotificationChannel.Publish(
		"",
		"send_notification",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		},
	)
	return err
}

func PublishSetRedis(key string, value string) error {
	body := key + "|" + value

	zap.L().Info("Publishing set_redis message to RabbitMQ: " + body)

	err := SetRedisChannel.Publish(
		"",
		"set_redis",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		},
	)
	return err
}

func PublishDeleteRedis(key string) error {
	body := key

	zap.L().Info("Publishing delete_redis message to RabbitMQ: " + body)

	err := SetRedisChannel.Publish(
		"",
		"delete_redis",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		},
	)
	return err
}

func PublishSendMention(username string, mentioner string) error {
	body := username + "|" + mentioner

	zap.L().Info("Publishing mention_user message to RabbitMQ: " + body)

	err := SendMentionChannel.Publish(
		"",
		"mention_user",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		},
	)
	return err
}

func PublishSendRepost(senderId string, threadId string, threadOwnerId string) error {
	body := senderId + "|" + threadId + "|" + threadOwnerId

	zap.L().Info("Publishing send_repost message to RabbitMQ: " + body)

	err := SendRepostChannel.Publish(
		"",
		"send_repost",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		},
	)
	return err
}
