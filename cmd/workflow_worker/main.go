package main

import (
	"log"

	"github.com/joho/godotenv"
	workflowworker "github.com/v-venes/friends-achievements-bot/internal/workflow_worker"
	"github.com/v-venes/friends-achievements-bot/pkg"
	steamclient "github.com/v-venes/friends-achievements-bot/pkg/steam_client"
	"go.temporal.io/sdk/client"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	env := pkg.GetEnvVars()
	log.Printf("env %s", env.TemporalHost)

	steamClient := steamclient.NewSteamClient(steamclient.NewSteamClientParams{
		SteamKey: env.SteamKey,
	})

	temporalClient, err := client.Dial(client.Options{
		HostPort: env.TemporalHost,
	})

	if err != nil {
		log.Fatalf("Erro ao conectar com Temporal: %s", err.Error())
	}

	workflowWorker := workflowworker.NewWorkflowWorker(workflowworker.NewWorkflowWorkerParams{
		Client:      temporalClient,
		SteamClient: steamClient,
	})

	workflowWorker.Run()
}
