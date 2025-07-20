package pkg

import (
	"log"

	"github.com/Netflix/go-env"
)

type Environment struct {
	DiscordBotToken string `env:"DISCORD_BOT_TOKEN"`
	DiscordGuildID  string `env:"DISCORD_GUILD_ID"`
	BrokerHost      string `env:"BROKER_HOST"`
	BrokerUsername  string `env:"BROKER_USERNAME"`
	BrokerPassword  string `env:"BROKER_PASSWORD"`
}

func GetEnvVars() *Environment {
	var config Environment
	_, err := env.UnmarshalFromEnviron(&config)
	if err != nil {
		log.Fatal(err)
	}

	return &config
}
