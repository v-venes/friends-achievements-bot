package workflows

import (
	"time"

	"github.com/v-venes/friends-achievements-bot/internal/workflow_worker/activities"
	"github.com/v-venes/friends-achievements-bot/pkg/repository"
	steamclient "github.com/v-venes/friends-achievements-bot/pkg/steam_client"
	"go.temporal.io/sdk/workflow"
)

func ExtractNewPlayerDataWorkflow(ctx workflow.Context, steamID string) error {
	logger := workflow.GetLogger(ctx)
	logger.Info("ExtractRecentlyPlayedGamesWorkflow started for", "steamID", steamID)

	activityOptions := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
	}
	ctx = workflow.WithActivityOptions(ctx, activityOptions)

	var recentlyPlayedGames steamclient.RecentlyPlayedGames
	err := workflow.ExecuteActivity(ctx, (*activities.SteamActivities).GetRecentlyPlayedGames, steamID).Get(ctx, &recentlyPlayedGames)

	if err != nil {
		return err
	}

	if len(recentlyPlayedGames.Games) == 0 {
		logger.Info("No recently played games found, finishing workflow")
		return nil
	}

	logger.Info("Ready to save found games")

	var gamesAchievements map[int]*steamclient.GameStats
	err = workflow.ExecuteActivity(ctx, (*activities.SteamActivities).GetAchievementsForGames, steamID, recentlyPlayedGames).Get(ctx, &gamesAchievements)
	if err != nil {
		return err
	}

	err = workflow.ExecuteActivity(ctx, (*activities.PlayerActivities).SavePlayerGamesAchievements, steamID, recentlyPlayedGames, gamesAchievements).Get(ctx, nil)
	if err != nil {
		return err
	}

	// err = workflow.ExecuteActivity(ctx, (*activities.QueueActivities).SendMessageToQueue, broker.ExtractGames, "").Get(ctx, nil)
	//

	return nil
}

func ExtractGameDataWorkflow(ctx workflow.Context, appID int) error {
	logger := workflow.GetLogger(ctx)
	logger.Info("ExtractGamesdataWorkflow started for", "appID", appID)

	activityOptions := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
	}
	ctx = workflow.WithActivityOptions(ctx, activityOptions)

	var existingGameDetails *repository.Game
	err := workflow.ExecuteActivity(ctx, (*activities.PlayerActivities).GetGameDataByAppID, appID).Get(ctx, &existingGameDetails)

	if err != nil {
		return err
	}

	if existingGameDetails != nil {
		return nil
	}

	var gameDetails steamclient.GameDetailsData
	err = workflow.ExecuteActivity(ctx, (*activities.SteamActivities).GetGameDetails, appID).Get(ctx, &gameDetails)

	if err != nil {
		return err
	}

	var gameAchievements steamclient.AllGameAchievements
	err = workflow.ExecuteActivity(ctx, (*activities.SteamActivities).GetAllGameAchievements, appID).Get(ctx, &gameAchievements)

	if err != nil {
		return err
	}

	err = workflow.ExecuteActivity(ctx, (*activities.PlayerActivities).SaveGameData, &gameDetails, &gameAchievements).Get(ctx, nil)

	if err != nil {
		return err
	}

	return nil
}
