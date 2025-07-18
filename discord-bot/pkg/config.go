package pkg

import (
	"log"

	"github.com/Netflix/go-env"
)

type Environment struct {
	DiscordBotToken string `env:"DISCORD_BOT_TOKEN"`
}

func GetEnvVars() *Environment {
	var config Environment
	_, err := env.UnmarshalFromEnviron(&config)
	if err != nil {
		log.Fatal(err)
	}

	return &config
}
