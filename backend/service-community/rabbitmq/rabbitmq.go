package rabbitmq

import (
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

var SendNotificationChannel *amqp.Channel

func InitRabbitMQ() {
	conn, err := amqp.Dial(os.Getenv("RABBITMQ_URL"))
	if err != nil {
		zap.L().Fatal("Failed to connect to RabbitMQ: " + err.Error())
	} else {
		zap.L().Info("Connected to RabbitMQ successfully.")
	}
	InitSendNotificationChannel(conn)
}

func InitSendNotificationChannel(conn *amqp.Connection) {
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
