package main

import (
	"log"

	"github.com/joho/godotenv"
	discordbot "github.com/v-venes/friends-achievements-bot/internal/discord-bot"
	"github.com/v-venes/friends-achievements-bot/pkg"
	"github.com/v-venes/friends-achievements-bot/pkg/broker"
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

	bot, err := discordbot.NewBot(discordbot.NewBotParams{
		DiscordToken:   env.DiscordBotToken,
		DisocrdGuildID: env.DiscordGuildID,
		Broker:         broker,
	})
	if err != nil {
		log.Fatalf("Erro ao criar bot %s", err.Error())
	}

	bot.Run()

	defer broker.Channel.Close()
	defer broker.Connection.Close()
}
