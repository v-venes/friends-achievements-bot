package queueworker

import (
	"context"
	"fmt"
	"sync"

	"github.com/rabbitmq/amqp091-go"
	"github.com/v-venes/friends-achievements-bot/internal/queue_worker/handlers"
	"github.com/v-venes/friends-achievements-bot/pkg/broker"
	"github.com/v-venes/friends-achievements-bot/pkg/repository"
	steamclient "github.com/v-venes/friends-achievements-bot/pkg/steam_client"
)

type QueueWorker struct {
	Broker           *broker.Broker
	Queues           map[string]func(amqp091.Delivery)
	SteamClient      *steamclient.SteamClient
	PlayerRepository *repository.PlayerRepository
}

type NewQueueWorkerParams struct {
	Broker           *broker.Broker
	SteamClient      *steamclient.SteamClient
	PlayerRepository *repository.PlayerRepository
}

func NewQueueWorker(params NewQueueWorkerParams) *QueueWorker {
	queueworker := &QueueWorker{
		Broker:           params.Broker,
		Queues:           map[string]func(amqp091.Delivery){},
		SteamClient:      params.SteamClient,
		PlayerRepository: params.PlayerRepository,
	}

	queueworker.registerHandlers()

	return queueworker
}

func (w *QueueWorker) registerHandlers() {
	w.Queues[broker.BrokerQueues[broker.NewSteamId]] = handlers.NewSteamIDHandler(handlers.NewSteamIDHandlerParams{
		PlayerRepository: w.PlayerRepository,
		SteamClient:      w.SteamClient,
		Broker:           w.Broker,
	})
}

func (w *QueueWorker) Run(ctx context.Context) error {
	var wg sync.WaitGroup

	for queue_name, handler := range w.Queues {
		msgs, err := w.Broker.ReceiveMessages(queue_name)
		if err != nil {
			return fmt.Errorf("failed to consume %s: %w", queue_name, err)
		}

		wg.Add(1)
		go func(m <-chan amqp091.Delivery, h func(amqp091.Delivery)) {
			defer wg.Done()
			for d := range m {
				h(d)
			}
		}(msgs, handler)
	}

	<-ctx.Done()

	w.Broker.Channel.Close()
	w.Broker.Connection.Close()
	wg.Wait()
	return nil
}
