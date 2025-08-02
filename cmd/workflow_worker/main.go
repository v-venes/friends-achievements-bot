package main

import (
	"context"
	"fmt"
	"log"

	"github.com/joho/godotenv"
	workflowworker "github.com/v-venes/friends-achievements-bot/internal/workflow_worker"
	"github.com/v-venes/friends-achievements-bot/pkg"
	"github.com/v-venes/friends-achievements-bot/pkg/repository"
	steamclient "github.com/v-venes/friends-achievements-bot/pkg/steam_client"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
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

	mongoUri := fmt.Sprintf("mongodb://%s/", env.MongoHost)
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(mongoUri).SetServerAPIOptions(serverAPI)
	mongoClient, err := mongo.Connect(opts)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err = mongoClient.Disconnect(context.TODO()); err != nil {
			log.Fatalf("Erro ao conectar com mongo: %v", err)
		}
	}()

	playerRepository := &repository.PlayerRepository{
		MongoClient: mongoClient,
	}

	temporalClient, err := client.Dial(client.Options{
		HostPort: env.TemporalHost,
	})

	if err != nil {
		log.Fatalf("Erro ao conectar com Temporal: %s", err.Error())
	}

	workflowWorker := workflowworker.NewWorkflowWorker(workflowworker.NewWorkflowWorkerParams{
		Client:           temporalClient,
		SteamClient:      steamClient,
		PlayerRepository: playerRepository,
	})

	workflowWorker.Run()
}
