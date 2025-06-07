package rabbitmq

import (
	"os"
	"strings"

	"github.com/nathabuddhi/ay-com/backend/service-user/handlers"
	"github.com/nathabuddhi/ay-com/backend/service-user/models"
	rabbitmq "github.com/nathabuddhi/ay-com/backend/service-user/rabbitmqsend"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

var ReceiveMentionChannel *amqp.Channel
var ReceiveRepostChannel *amqp.Channel

var userHandler *handlers.Handlers
var Forever = make(chan bool)

func InitHandler(newHandler *handlers.Handlers) {
	conn, err := amqp.Dial(os.Getenv("RABBITMQ_URL"))
	if err != nil {
		zap.L().Fatal("Failed to connect to RabbitMQ: " + err.Error())
	} else {
		zap.L().Info("Connected to RabbitMQ successfully.")
	}
	userHandler = newHandler
	InitReceiverChannel(conn)
	InitRepostChannel(conn)

	StartConsumingMentions()
	StartConsumingReposts()
	<-Forever
}

func InitReceiverChannel(conn *amqp.Connection) {
	ch, err := conn.Channel()
	if err != nil {
		zap.L().Fatal("Failed to open a channel: " + err.Error())
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

	ReceiveMentionChannel = ch
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

	ReceiveRepostChannel = ch
}

func StartConsumingMentions() {
	msgs, err := ReceiveMentionChannel.Consume(
		"mention_user",
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		zap.L().Fatal("Failed to register mention_user consumer: " + err.Error())
	} else {
		zap.L().Info("RabbitMQ mention_user consumer registered successfully")
	}

	go func() {
		for d := range msgs {
			zap.L().Info("Received a message: " + string(d.Body))
			data := strings.Split(string(d.Body), "|")
			username := data[0]
			mentionerId := data[1]
			if username == "" || mentionerId == "" {
				zap.L().Error("Received an empty username or mentioner in mention_user queue")
				continue
			} else {
				var user models.User
				var mentioner models.User
				if err := userHandler.DB.Where("username = ?", username).First(&user).Error; err != nil {
					zap.L().Error("Failed to find user by username: " + username + " - " + err.Error())
				} else {
					if err := userHandler.DB.Where("user_id = ?", mentionerId).First(&mentioner).Error; err != nil {
						zap.L().Error("Failed to find user by user_id: " + mentionerId + " - " + err.Error())
					}

					rabbitmq.PublishSendNotification("mention", user.UserId, user.Email, mentioner.Username+" Mentioned you in a thread.", "You have been mentioned by <a href='/profile/"+mentioner.Username+"'>"+mentioner.Username+"</a> in a post.", mentioner.Name)
				}
			}
		}
	}()

	zap.L().Info("mention_user listener running.")
}

func StartConsumingReposts() {
	msgs, err := ReceiveRepostChannel.Consume(
		"send_repost",
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		zap.L().Fatal("Failed to register send_repost consumer: " + err.Error())
	} else {
		zap.L().Info("RabbitMQ send_repost consumer registered successfully")
	}

	go func() {
		for d := range msgs {
			zap.L().Info("Received a message: " + string(d.Body))
			data := strings.Split(string(d.Body), "|")
			senderId := data[0]
			threadId := data[1]
			threadOwnerId := data[2]
			if senderId == "" || threadId == "" || threadOwnerId == "" {
				zap.L().Error("Received an empty senderId, threadId or threadOwnerId in send_repost queue")
				continue
			} else {
				var user models.User
				var threadOwner models.User
				if err := userHandler.DB.Where("user_id = ?", senderId).First(&user).Error; err != nil {
					zap.L().Error("Failed to find user by user_id: " + senderId + " - " + err.Error())
				} else {
					if err := userHandler.DB.Where("user_id = ?", threadOwnerId).First(&threadOwner).Error; err != nil {
						zap.L().Error("Failed to find user by user_id: " + threadOwnerId + " - " + err.Error())
					}

					rabbitmq.PublishSendNotification("repost", threadOwner.UserId, threadOwner.Email, user.Username+" Reposted your thread.", "<a href='/profile/"+user.Username+"'>"+user.Username+"</a> Reposted <a href='/thread/"+threadId+"'>Your Thread</a>", user.Name)
				}
			}
		}
	}()

	zap.L().Info("mention_user listener running.")
}
