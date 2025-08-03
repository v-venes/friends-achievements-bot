package repository

import (
	"context"

	steamclient "github.com/v-venes/friends-achievements-bot/pkg/steam_client"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

const GAMES_COLLECTION = "games"

type Game struct {
	AppID            int               `bson:"app_id"`
	AppName          string            `bson:"app_name"`
	ShortDescription string            `bson:"short_description"`
	HeaderImage      string            `bson:"header_image"`
	Achievements     []GameAchievement `bson:"achievements"`
}

type GameAchievement struct {
	AchievementID string `bson:"achievement_id"`
	DisplayName   string `bson:"display_name"`
	Description   string `bson:"description"`
	Icon          string `bson:"icon"`
}

func NewGameDataFromSteam(gameDetails *steamclient.GameDetailsData, gameAchievements *steamclient.AllGameAchievements) *Game {
	var achievements []GameAchievement

	for _, achievement := range gameAchievements.AvaiableGameStats.Achievements {
		achievements = append(achievements, GameAchievement{
			AchievementID: achievement.Name,
			DisplayName:   achievement.DisplayName,
			Description:   achievement.Description,
			Icon:          achievement.Icon,
		})
	}

	return &Game{
		AppID:            gameDetails.AppID,
		AppName:          gameDetails.Name,
		ShortDescription: gameDetails.ShortDescription,
		HeaderImage:      gameDetails.HeaderImage,
		Achievements:     achievements,
	}
}

func (pr *PlayerRepository) GetGameByAppID(appID int) (*Game, error) {
	gamesCollection := pr.MongoClient.Database(DEFAULT_DATABASE).Collection(GAMES_COLLECTION)
	filter := bson.M{
		"app_id": appID,
	}
	var game Game
	err := gamesCollection.FindOne(context.TODO(), filter).Decode(&game)

	if err != nil {
		if err != mongo.ErrNoDocuments {
			return nil, err
		}
		return nil, nil
	}

	return &game, nil
}

func (pr *PlayerRepository) SaveGameData(game *Game) error {
	gamesCollection := pr.MongoClient.Database(DEFAULT_DATABASE).Collection(GAMES_COLLECTION)
	_, err := gamesCollection.InsertOne(context.TODO(), game)
	return err
}
