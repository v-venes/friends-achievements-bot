package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/v-venes/friends-achievements-bot/pkg"
	"github.com/v-venes/friends-achievements-bot/pkg/repository"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

// TODO: Rodar criação dos indices no docker-compose
func main() {
	env := pkg.GetEnvVars()

	mongoUri := fmt.Sprintf("mongodb://%s/", env.MongoHost)
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(mongoUri).SetServerAPIOptions(serverAPI)
	mongoClient, err := mongo.Connect(opts)
	if err != nil {
		panic(err)
	}

	database := mongoClient.Database(repository.DEFAULT_DATABASE)

	err = createIndexes(database)
	if err != nil {
		log.Fatal("Erro ao criar índices:", err)
	}

	defer func() {
		if err = mongoClient.Disconnect(context.TODO()); err != nil {
			log.Fatalf("Erro ao conectar com mongo: %v", err)
		}
	}()
}

func createIndexes(database *mongo.Database) error {
	playerCollection := database.Collection(repository.PLAYER_COLLECTION)
	playerIndex := mongo.IndexModel{
		Keys: bson.D{
			{Key: "player_id", Value: 1},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := playerCollection.Indexes().CreateOne(ctx, playerIndex)
	if err != nil {
		return err
	}

	playerAchievementsCollection := database.Collection(repository.PLAYER_GAMES_ACHIEVEMENTS_COLLECTION)
	playersAchievementsIndex := mongo.IndexModel{
		Keys: bson.D{
			{Key: "player_id", Value: 1},
			{Key: "appid", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	}

	_, err = playerAchievementsCollection.Indexes().CreateOne(ctx, playersAchievementsIndex)
	if err != nil {
		return err
	}
	return nil
}
