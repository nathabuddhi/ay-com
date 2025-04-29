package rabbitmq

import (
	"strings"

	"github.com/nathabuddhi/ay-com/backend/util-email/email"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

func StartConsuming() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		zap.L().Fatal("Failed to connect to RabbitMQ: " + err.Error())
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		zap.L().Fatal("Failed to open a channel: " + err.Error())
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"send_email",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		zap.L().Fatal("Failed to declare send_email queue: " + err.Error())
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		zap.L().Fatal("Failed to register send_email consumer: " + err.Error())
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			zap.L().Info("Received a message: " + string(d.Body))

			data := strings.Split(string(d.Body), "|")
			if len(data) != 3 {
				zap.L().Error("Invalid message format.")
				continue
			}

			email := data[0]
			subject := data[1]
			body := data[2]

			sendEmail(email, subject, body)
		}
	}()

	zap.L().Info(" [*] Waiting for messages...")
	<-forever
}

func sendEmail(toEmail string, subject string, body string) {
	zap.L().Info("Sending verification email to " + toEmail + " with subject " + subject + ".")
	email.SendEmail(toEmail, subject, body)
}
