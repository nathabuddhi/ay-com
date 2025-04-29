package rabbitmq

import (
	"github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

var Channel *amqp091.Channel

func InitRabbitMQ() {
	conn, err := amqp091.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		zap.L().Fatal("Failed to connect to RabbitMQ: " + err.Error())
	}

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

	Channel = ch
}

func PublishEmail(email string, subject string, body string) error {
	full_body := email + "|" + subject + "|" + body

	zap.L().Info("Publishing notification message to RabbitMQ: " + full_body)

	err := Channel.Publish(
		"",
		"send_email",
		false,
		false,
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        []byte(full_body),
		},
	)
	return err
}
