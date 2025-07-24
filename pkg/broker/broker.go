package broker

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/rabbitmq/amqp091-go"
)

type BrokerQueueEnum int

const (
	NewSteamId BrokerQueueEnum = iota
	NewAchievement
	FeedbackMessages
)

const APP_BROKER_EXCHANGE = "ACHIEVEMENTS-BOT"

type NewBroketParams struct {
	Username string
	Password string
	Host     string
}

type Broker struct {
	Connection *amqp091.Connection
	Channel    *amqp091.Channel
	Queues     map[BrokerQueueEnum]*amqp091.Queue
}

type SendMessageParams struct {
	Queue   BrokerQueueEnum
	Message any
}

var BrokerQueues = map[BrokerQueueEnum]string{
	NewSteamId:       "NEW_STEAM_ID",
	NewAchievement:   "NEW_ACHIEVEMENT",
	FeedbackMessages: "FEEDBACK_MESSAGES",
}

func NewBroker(params NewBroketParams) (*Broker, error) {
	conn, err := amqp091.Dial(fmt.Sprintf("amqp://%s:%s@%s", params.Username, params.Password, params.Host))
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	queues := make(map[BrokerQueueEnum]*amqp091.Queue)
	for key, queueName := range BrokerQueues {
		q, err := ch.QueueDeclare(queueName, true, false, false, false, nil)
		if err != nil {
			return nil, err
		}

		queues[key] = &q
	}

	return &Broker{
		Connection: conn,
		Channel:    ch,
		Queues:     queues,
	}, nil
}

func (b *Broker) SendMessage(params SendMessageParams) error {
	body, err := json.Marshal(params.Message)
	if err != nil {
		return nil
	}

	err = b.Channel.Publish("", b.Queues[params.Queue].Name, false, false, amqp091.Publishing{
		ContentType:  "application/json",
		Body:         body,
		DeliveryMode: 2,
	})
	if err != nil {
		log.Printf("Erro ao enviar mensagem [%s]: %s", b.Queues[params.Queue].Name, err.Error())
		return err
	}

	log.Printf("Mensagem enviada para queue %s", b.Queues[params.Queue].Name)
	return nil
}

func (b *Broker) ReceiveMessages(queue string) (<-chan amqp091.Delivery, error) {
	msgs, err := b.Channel.Consume(queue, "", true, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	return msgs, nil
}
