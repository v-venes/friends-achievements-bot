package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/v-venes/friends-achievements-bot/cmd/workflow_worker/workflows"
	"github.com/v-venes/friends-achievements-bot/pkg"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
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

	temporalClient, err := client.Dial(client.Options{
		HostPort: env.TemporalHost,
	})

	if err != nil {
		log.Fatalf("Erro ao conectar com Temporal: %s", err.Error())
	}

	temporalWorker := worker.New(temporalClient, "steam-achievements", worker.Options{})

	worklfowOptions := workflow.RegisterOptions{
		Name: "ExtractPlayerGames",
	}
	temporalWorker.RegisterWorkflowWithOptions(workflows.ExtractPlayerGames, worklfowOptions)

	err = temporalWorker.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalf("Erro ao iniciar worker: %s", err.Error())
	}

	// Extrair jogos recentes (assim que adicionar o steamID)

}
