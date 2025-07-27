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
	SteamKey        string `env:"STEAM_KEY"`
	MongoHost       string `env:"MONGO_HOST"`
	MongoUsername   string `env:"MONGO_USERNAME"`
	MongoPassword   string `env:"MONGO_PASSWORD"`
	TemporalHost    string `env:"TEMPORAL_HOST"`
}

func GetEnvVars() *Environment {
	var config Environment
	_, err := env.UnmarshalFromEnviron(&config)
	if err != nil {
		log.Fatal(err)
	}

	return &config
}
