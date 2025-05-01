package rabbitmq

import (
	"os"

	"github.com/nathabuddhi/ay-com/backend/service-notification/handlers"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

var EmailChannel *amqp.Channel
var ReceiverChannel *amqp.Channel
var Handler *handlers.Handler

func InitReceiverChannel(conn *amqp.Connection) {
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

	ReceiverChannel = ch
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

func InitRabbitMQ(h *handlers.Handler) {
	conn, err := amqp.Dial(os.Getenv("RABBITMQ_URL"))
	if err != nil {
		zap.L().Fatal("Failed to connect to RabbitMQ: " + err.Error())
	} else {
		zap.L().Info("Connected to RabbitMQ successfully.")
	}
	InitEmailChannel(conn)
	InitReceiverChannel(conn)
	Handler = h
	
	StartConsumingNotifications()
}
