package rabbitmq

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/microphoneabuser/balance-service/models"
	"github.com/spf13/viper"
	"github.com/streadway/amqp"
)

type RabbitConfig struct {
	User     string
	Password string
	Host     string
	Port     string
}

func NewRabbitMQConn(config *RabbitConfig) (*amqp.Connection, error) {
	connAddr := fmt.Sprintf(
		"amqp://%s:%s@%s:%s/",
		config.User,
		config.Password,
		config.Host,
		config.Port,
	)
	return amqp.Dial(connAddr)
}

var ch *amqp.Channel = &amqp.Channel{}
var q *amqp.Queue = &amqp.Queue{}

func SetQueue(conn amqp.Connection) error {
	var err error
	ch, err = conn.Channel()
	if err != nil {
		return errors.New("Failed to open channel")
	}

	*q, err = ch.QueueDeclare(
		viper.GetString("rabbitmq.Queue"), // name
		false,                             // durable
		false,                             // delete when unused
		false,                             // exclusive
		false,                             // no-wait
		nil,                               // arguments
	)
	if err != nil {
		return errors.New("Failed to declare queue")
	}
	return nil
}

func CloseChannel() {
	ch.Close()
}

type SMS struct {
	AccountId int    `json:"account_id"`
	Message   string `json:"message"`
}

func PublishAccrualDebiting(data models.AccountInput, balance int, isAccrual bool) {
	finalAmount := toNormal(data.Amount)
	finalBalance := toNormal(balance)

	var sms SMS

	if isAccrual {
		sms = SMS{
			AccountId: data.Id,
			Message:   fmt.Sprintf("Счет-%d Зачисление %.2fр Баланс: %.2fр", data.Id, finalAmount, finalBalance),
		}
	} else {
		sms = SMS{
			AccountId: data.Id,
			Message:   fmt.Sprintf("Счет-%d Списание %.2fр Баланс: %.2fр", data.Id, finalAmount, finalBalance),
		}
	}

	publish(sms)
}

func PublishTransfer(data models.TransactionInput, senderBalance int, recipientBalance int) {
	finalAmount := toNormal(data.Amount)
	finalSenderBalance := toNormal(senderBalance)

	smsToSender := SMS{
		AccountId: data.SenderId,
		Message:   fmt.Sprintf("Счет-%d Исходящий перевод %.2fр Баланс: %.2fр", data.SenderId, finalAmount, finalSenderBalance),
	}

	publish(smsToSender)

	finalRecipientBalance := toNormal(recipientBalance)

	smsToRecipient := SMS{
		AccountId: data.RecipientId,
		Message: fmt.Sprintf("Счет-%d Входящий перевод %.2fр от Счет-%d Баланс: %.2fр Сообщение:\"%s\"",
			data.RecipientId, finalAmount, data.SenderId, finalRecipientBalance, data.Description),
	}

	publish(smsToRecipient)
}

func publish(data interface{}) {
	body, err := json.Marshal(data)
	if err != nil {
		log.Printf("RabbitMQ: Failed to encode json object: %s", err.Error())
		return
	}

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		log.Printf("RabbitMQ: Failed to publish: %s", err.Error())
		return
	}
}

func toNormal(amount int) float64 {
	return float64(amount) / 100
}
