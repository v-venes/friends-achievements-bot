package activities

import (
	"github.com/v-venes/friends-achievements-bot/pkg/repository"
	steamclient "github.com/v-venes/friends-achievements-bot/pkg/steam_client"
)

type PlayerActivities struct {
	PlayerRepository *repository.PlayerRepository
}

func (p *PlayerActivities) SaveRecentlyPlayedGames(recentlyPlayedGames steamclient.RecentlyPlayedGames) error {
	p.PlayerRepository.CreatePlayer()

	return nil
}
