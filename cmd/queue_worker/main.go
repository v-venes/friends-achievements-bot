package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	queueworker "github.com/v-venes/friends-achievements-bot/internal/queue_worker"
	"github.com/v-venes/friends-achievements-bot/pkg"
	"github.com/v-venes/friends-achievements-bot/pkg/broker"
	"github.com/v-venes/friends-achievements-bot/pkg/service"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	env := pkg.GetEnvVars()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signalChan
		cancel()
	}()

	newBroker, err := broker.NewBroker(broker.NewBroketParams{
		Username: env.BrokerUsername,
		Password: env.BrokerPassword,
		Host:     env.BrokerHost,
	})
	if err != nil {
		log.Fatalf("Erro ao conectar com broker %s", err.Error())
	}

	steamClient := service.NewSteamClient(service.NewSteamClientParams{
		SteamKey: env.SteamKey,
	})

	worker := queueworker.NewQueueWorker(queueworker.NewQueueWorkerParams{
		Broker:      newBroker,
		SteamClient: steamClient,
	})

	if err := worker.Run(ctx); err != nil {
		log.Fatalf("Erro ao iniciar worker: %v", err)
	}
}
