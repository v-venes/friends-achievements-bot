package activities

import (
	"context"

	steamclient "github.com/v-venes/friends-achievements-bot/pkg/steam_client"
	"go.temporal.io/sdk/activity"
)

type SteamActivities struct {
	Client *steamclient.SteamClient
}

func (s *SteamActivities) ExtractRecentlyPlayedGames(ctx context.Context, steamID string) (*steamclient.RecentlyPlayedGames, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("ExtractRecentlyPlayedGames started for", "steamID", steamID)

	recentlyPlayedGames, err := s.Client.GetRecentlyPlayedGames(steamID)

	if err != nil {
		logger.Error("Get Recently Playerd Games", "Error", err.Error())
		return nil, err
	}

	return recentlyPlayedGames, nil
}
