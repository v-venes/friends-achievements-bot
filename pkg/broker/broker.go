package broker

import (
	"encoding/json"
	"fmt"

	"github.com/rabbitmq/amqp091-go"
)

type BrokerQueueEnum int

const (
	NewSteamId BrokerQueueEnum = iota
	NewAchievement
)

type NewBroketParams struct {
	Username string
	Password string
	Host     string
}

type Broker struct {
	Channel *amqp091.Channel
	Queues  map[BrokerQueueEnum]amqp091.Queue
}

type SendMessageParams struct {
	Queue   BrokerQueueEnum
	Message any
}

var BrokerQueues = map[BrokerQueueEnum]string{
	NewSteamId:     "NEW_STEAM_ID",
	NewAchievement: "NEW_ACHIEVEMENT",
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

	queues := make(map[BrokerQueueEnum]amqp091.Queue)
	for key, v := range BrokerQueues {
		q, err := ch.QueueDeclare(v, false, false, false, false, nil)
		if err != nil {
			return nil, err
		}

		queues[key] = q
	}

	defer ch.Close()
	defer conn.Close()
	return &Broker{
		Channel: ch,
		Queues:  queues,
	}, nil
}

func (b *Broker) SendMessage(params SendMessageParams) error {
	body, err := json.Marshal(params.Message)
	if err != nil {
		return nil
	}

	b.Channel.Publish("", b.Queues[params.Queue].Name, false, false, amqp091.Publishing{
		ContentType: "application/json",
		Body:        body,
	})

	return nil
}

func (b *Broker) ReceiveMessages(queue string) (<-chan amqp091.Delivery, error) {
	msgs, err := b.Channel.Consume(queue, "", true, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	return msgs, nil
}
