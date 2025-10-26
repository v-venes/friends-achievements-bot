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

func (p *PlayerActivities) GetNotFoundGames(appIDs []int) (*[]int, error) {
	games, err := p.PlayerRepository.FindGamesByAppID(appIDs)

	if err != nil {
		return nil, err
	}

	gamesSet := make(map[int]*repository.Game)
	for _, game := range *games {
		gamesSet[game.AppID] = &game
	}

	var notFoundIDs []int

	for _, appID := range appIDs {
		_, ok := gamesSet[appID]
		if !ok {
			notFoundIDs = append(notFoundIDs, appID)
		}
	}

	return &notFoundIDs, nil
}

func (p *PlayerActivities) GetAllPlayers() (*[]repository.PlayerModel, error) {
	players, err := p.PlayerRepository.GetAll()

	if err != nil {
		return nil, err
	}

	return players, nil
}
