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
