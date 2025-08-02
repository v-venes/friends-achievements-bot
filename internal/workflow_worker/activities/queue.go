package activities

import "github.com/v-venes/friends-achievements-bot/pkg/broker"

type QueueActivities struct {
	broker *broker.Broker
}

func (qa *QueueActivities) SendMessageToQueue(queue broker.BrokerQueueEnum, message any) {
	qa.broker.SendMessage(broker.SendMessageParams{
		Queue:   queue,
		Message: message,
	})

}
