package repository

import (
	"context"
	"time"

	steamclient "github.com/v-venes/friends-achievements-bot/pkg/steam_client"
	"go.mongodb.org/mongo-driver/v2/bson"
)

const PLAYER_GAMES_ACHIEVEMENTS_COLLECTION = "players_achievements"

type PlayerGameAchievements struct {
	AppID        int                     `bson:"app_id"`
	PlayerID     string                  `bson:"player_id"`
	AppName      string                  `bson:"app_name"`
	LastChecked  time.Time               `bson:"last_checked"`
	Achievements []PlayerGameAchievement `bson:"achievements"`
}

type PlayerGameAchievement struct {
	AchievementID string    `bson:"achievement_id"`
	Achieved      bool      `bson:"achieved"`
	UnclockedAt   time.Time `bson:"unlocked_at"`
}

func NewPlayerGameAchievementsFromSteam(playerID string, recentlyPlayedGames steamclient.RecentlyPlayedGames, gamesAchievements map[int]*steamclient.GameStats) []*PlayerGameAchievements {
	var playerGamesAchievements []*PlayerGameAchievements

	for _, game := range recentlyPlayedGames.Games {
		//TODO: Melhorar o parse dos achievements
		achievements := []PlayerGameAchievement{}

		_, ok := gamesAchievements[game.AppID]
		if !ok {
			continue
		}

		for _, achievement := range gamesAchievements[game.AppID].Achievements {
			achievements = append(achievements, PlayerGameAchievement{
				AchievementID: achievement.Name,
				Achieved:      achievement.Achieved != 0,
				UnclockedAt:   time.Now(),
			})

		}

		gameAchievements := &PlayerGameAchievements{
			AppID:        game.AppID,
			PlayerID:     playerID,
			AppName:      game.Name,
			LastChecked:  time.Now(),
			Achievements: achievements,
		}
		playerGamesAchievements = append(playerGamesAchievements, gameAchievements)
	}

	return playerGamesAchievements
}

func (pr *PlayerRepository) FindManyByPlayerAndGames(playerID string, appIDs []string) (map[int]*PlayerGameAchievements, error) {
	playerGameAchievementsCollection := pr.MongoClient.Database(DEFAULT_DATABASE).Collection(PLAYER_GAMES_ACHIEVEMENTS_COLLECTION)
	filter := bson.M{
		"player_id": playerID,
		"app_id":    bson.M{"$in": appIDs},
	}

	cursor, err := playerGameAchievementsCollection.Find(context.TODO(), filter)

	if err != nil {
		return nil, err
	}

	var results []PlayerGameAchievements
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}

	achievementsMap := map[int]*PlayerGameAchievements{}

	for _, result := range results {
		achievementsMap[result.AppID] = &result
	}

	return achievementsMap, nil
}

func (pr *PlayerRepository) CreatePlayerGamesAchievements(gamesAchievements []*PlayerGameAchievements) error {
	playerAchievementsCollection := pr.MongoClient.Database(DEFAULT_DATABASE).Collection(PLAYER_GAMES_ACHIEVEMENTS_COLLECTION)
	_, err := playerAchievementsCollection.InsertMany(context.TODO(), gamesAchievements)
	return err
}
