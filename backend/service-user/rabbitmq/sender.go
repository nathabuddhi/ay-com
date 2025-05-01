package rabbitmq

import (
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

var EmailChannel *amqp.Channel
var SetRedisChannel *amqp.Channel
var DeleteRedisChannel *amqp.Channel
var SendNotificationChannel *amqp.Channel

func InitNotificationChannel(conn *amqp.Connection) {
	ch, err := conn.Channel()
	if err != nil {
		zap.L().Fatal("Failed to open a channel: " + err.Error())
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

	SendNotificationChannel = ch
}

func InitEmailChannel(conn *amqp.Connection) {
	ch, err := conn.Channel()
	if err != nil {
		zap.L().Fatal("Failed to open a channel: " + err.Error())
	}

	_, err = ch.QueueDeclare(
		"send_email",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		zap.L().Fatal("Failed to declare send_email queue: " + err.Error())
	} else {
		zap.L().Info("RabbitMQ send_email queue declared successfully")
	}

	EmailChannel = ch
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

func InitRabbitMQ() {
	conn, err := amqp.Dial(os.Getenv("RABBITMQ_URL"))
	if err != nil {
		zap.L().Fatal("Failed to connect to RabbitMQ: " + err.Error())
	} else {
		zap.L().Info("Connected to RabbitMQ successfully.")
	}
	InitEmailChannel(conn)
	InitSetRedisChannel(conn)
	InitDeleteRedisChannel(conn)
	InitNotificationChannel(conn)
}

func PublishEmail(email string, subject string, body string) error {
	full_body := email + "|" + subject + "|" + body

	zap.L().Info("Publishing notification message to RabbitMQ: " + full_body)

	err := EmailChannel.Publish(
		"",
		"send_email",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(full_body),
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

func PublishSendNotification(notif_type string, user_id string, email string, title string, content string, from string) error {
	body := notif_type + "|" + user_id + "|" + email + "|" + title + "|" + content + "|" + from

	zap.L().Info("Publishing send_notification message to RabbitMQ: " + body)

	err := SendNotificationChannel.Publish(
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
