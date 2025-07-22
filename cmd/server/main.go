package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/v-venes/friends-achievements-bot/internal/server"
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

	broker, err := broker.NewBroker(broker.NewBroketParams{
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

	server := server.NewServer(server.NewServerParams{
		Broker:      broker,
		SteamClient: steamClient,
	})

	err = server.Run()
	if err != nil {
		log.Fatalf("Erro ao iniciar servidor %s", err.Error())
	}
}
