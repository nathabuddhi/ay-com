package rabbitmq

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

func SendEmail(toEmail string, subject string, body string) {
	zap.L().Info("Sending email to " + toEmail + " with subject " + subject + ".")

	full_body := toEmail + "|" + subject + "|" + body

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
	if err != nil {
		zap.L().Error("Failed to publish send_email message to RabbitMQ: " + err.Error())
	}
}
