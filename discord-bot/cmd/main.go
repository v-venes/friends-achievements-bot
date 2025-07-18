package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/v-venes/friends-achievements-bot/discord-bot/internal"
	"github.com/v-venes/friends-achievements-bot/discord-bot/pkg"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	env := pkg.GetEnvVars()
	bot, err := internal.NewBot(env.DiscordBotToken)
	if err != nil {
		log.Fatalf("Error creating bot: %s", err.Error())
	}

	bot.StartBot()
}
