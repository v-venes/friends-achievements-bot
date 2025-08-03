package activities

import (
	"context"

	steamclient "github.com/v-venes/friends-achievements-bot/pkg/steam_client"
	"go.temporal.io/sdk/activity"
)

type SteamActivities struct {
	Client *steamclient.SteamClient
}

func (s *SteamActivities) GetRecentlyPlayedGames(ctx context.Context, steamID string) (*steamclient.RecentlyPlayedGames, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("ExtractRecentlyPlayedGames started for", "steamID", steamID)

	recentlyPlayedGames, err := s.Client.GetRecentlyPlayedGames(steamID)

	if err != nil {
		logger.Error("Get Recently Playerd Games", "Error", err.Error())
		return nil, err
	}

	return recentlyPlayedGames, nil
}

func (s *SteamActivities) GetAchievementsForGames(ctx context.Context, steamID string, recentlyPlayedGames steamclient.RecentlyPlayedGames) (map[int]steamclient.GameStats, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("ExtractRecentlyPlayedGames started for", "steamID", steamID)

	gamesAchievements := map[int]steamclient.GameStats{}

	for _, game := range recentlyPlayedGames.Games {
		gameStats, err := s.Client.GetGameStats(steamID, game.AppID)

		if gameStats.PlayerStats.Achievements == nil {
			logger.Info("Achievements not found", "AppID", game.AppID)
			continue
		}

		gamesAchievements[game.AppID] = gameStats.PlayerStats

		if err != nil {
			logger.Error("Get Recently Playerd Games", "Error", err.Error())
			return nil, err
		}
	}

	return gamesAchievements, nil
}

func (s *SteamActivities) GetGameDetails(ctx context.Context, appID int) (*steamclient.GameDetailsData, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("GetGameDetails started for", "appID", appID)

	gameDetails, err := s.Client.GetGameDetails(appID)

	if err != nil {
		logger.Error("Get Recently Playerd Games", "Error", err.Error())
		return nil, err
	}

	return gameDetails, nil
}

func (s *SteamActivities) GetAllGameAchievements(ctx context.Context, appID int) (*steamclient.AllGameAchievements, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("GetAllGameAchievements started for", "appID", appID)

	gameAchievements, err := s.Client.GetAllGameAchievements(appID)

	if err != nil {
		logger.Error("Get Recently Playerd Games", "Error", err.Error())
		return nil, err
	}

	return gameAchievements, nil
}
