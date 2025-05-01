package rabbitmq

import (
	"strings"

	"go.uber.org/zap"
)

func StartConsumingNotifications() {
	msgs, err := ReceiverChannel.Consume(
		"send_notification",
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		zap.L().Fatal("Failed to register send_notification consumer: " + err.Error())
	} else {
		zap.L().Info("RabbitMQ send_notification consumer registered successfully")

	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			zap.L().Info("Received a message: " + string(d.Body))

			data := strings.Split(string(d.Body), "|")
			if len(data) != 6 {
				zap.L().Error("Invalid message format.")
				continue
			}

			notif_type := data[0]
			if notif_type != "follow" && notif_type != "mention" && notif_type != "like" && notif_type != "community" && notif_type != "repost" && notif_type != "newsletter" && notif_type != "system" {
				zap.L().Warn("Invalid notification type: " + notif_type)
			} else {
				user_id := data[1]
				email := data[2]
				title := data[3]
				content := data[4]
				from := data[5]

				notifOn := Handler.CheckNotificationSettings(user_id, notif_type)

				if notifOn {
					zap.L().Info("Creating notification for user: " + user_id)
					Handler.CreateNotification(user_id, title, content, from)
					SendEmail(email, "AY.com - Notification: "+title, "You have received a new notification: "+content+"<br><br>From: "+from+"<br><br>Thank you for using AY.com! You may disable notifications at any time via your account settings.")
				}
			}
		}
	}()

	zap.L().Info("Notification Service RabbitMQ listener running.")
	<-forever
}
