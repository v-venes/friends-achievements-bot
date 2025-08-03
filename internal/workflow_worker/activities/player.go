package activities

import (
	"github.com/v-venes/friends-achievements-bot/pkg/repository"
	steamclient "github.com/v-venes/friends-achievements-bot/pkg/steam_client"
)

type PlayerActivities struct {
	PlayerRepository *repository.PlayerRepository
}

func (p *PlayerActivities) SavePlayerGamesAchievements(steamID string, recentlyPlayedGames steamclient.RecentlyPlayedGames, gamesAchievements map[int]*steamclient.GameStats) error {

	playerGamesAchievements := repository.NewPlayerGameAchievementsFromSteam(steamID, recentlyPlayedGames, gamesAchievements)

	err := p.PlayerRepository.CreatePlayerGamesAchievements(playerGamesAchievements)

	if err != nil {
		return err
	}

	return nil
}

func (p *PlayerActivities) SaveGameData(gameDetails *steamclient.GameDetailsData, gameAchievements *steamclient.AllGameAchievements) error {

	gameData := repository.NewGameDataFromSteam(gameDetails, gameAchievements)

	err := p.PlayerRepository.SaveGameData(gameData)

	if err != nil {
		return err
	}

	return nil
}

func (p *PlayerActivities) GetGameDataByAppID(appID int) (*repository.Game, error) {
	game, err := p.PlayerRepository.GetGameByAppID(appID)

	if err != nil {
		return nil, err
	}

	return game, nil
}
