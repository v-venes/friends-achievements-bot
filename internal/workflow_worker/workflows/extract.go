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

	// TODO: Enviar mensagem para iniciar extração
	// err = workflow.ExecuteActivity(ctx, (*activities.QueueActivities).SendMessageToQueue, broker.ExtractGames, "").Get(ctx, nil)

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

func DailyAchievementsVerificationWorkflow(ctx workflow.Context) error {
	logger := workflow.GetLogger(ctx)
	logger.Info("DailyAchievementsVerificationWorkflow started")

	activityOptions := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
	}
	ctx = workflow.WithActivityOptions(ctx, activityOptions)

	var players *[]repository.PlayerModel
	err := workflow.ExecuteActivity(ctx, (*activities.PlayerActivities).GetAllPlayers).Get(ctx, &players)

	if err != nil {
		return err
	}

	futures := []workflow.ChildWorkflowFuture{}

	for _, player := range *players {
		future := workflow.ExecuteChildWorkflow(ctx, FindNewAchievementsWorkflow, player.PlayerID)
		futures = append(futures, future)
	}

	for i, f := range futures {
		err := f.Get(ctx, nil)

		if err != nil {
			logger.Error("Child workflow failed", "playerID", (*players)[i].PlayerID, "error", err)
		}
	}

	return nil
}

func FindNewAchievementsWorkflow(ctx workflow.Context, playerID string) error {
	logger := workflow.GetLogger(ctx)
	logger.Info("FindPlayerNewAchievementWorkflow started for", "playerID", playerID)

	activityOptions := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
	}
	ctx = workflow.WithActivityOptions(ctx, activityOptions)

	var recentlyPlayedGames steamclient.RecentlyPlayedGames
	err := workflow.ExecuteActivity(ctx, (*activities.SteamActivities).GetRecentlyPlayedGames, playerID).Get(ctx, &recentlyPlayedGames)

	if err != nil {
		return err
	}

	if len(recentlyPlayedGames.Games) == 0 {
		logger.Info("No recently played games found, finishing workflow")
		return nil
	}

	var appIds []int
	for _, game := range recentlyPlayedGames.Games {
		appIds = append(appIds, game.AppID)
	}

	var gamesToExtract []int
	err = workflow.ExecuteActivity(ctx, (*activities.PlayerActivities).GetNotFoundGames, appIds).Get(ctx, &gamesToExtract)

	// TODO: Se dois workflows tentarem extrair o mesmo jogo?
	if len(gamesToExtract) > 0 {
		futures := []workflow.ChildWorkflowFuture{}

		for _, gameID := range gamesToExtract {
			future := workflow.ExecuteChildWorkflow(ctx, ExtractGameDataWorkflow, gameID)
			futures = append(futures, future)
		}

		for i, f := range futures {
			err := f.Get(ctx, nil)

			if err != nil {
				logger.Error("Child workflow failed", "gameID", gamesToExtract[i], "error", err)
			}
		}
	}

	var gamesAchievements map[int]*steamclient.GameStats
	err = workflow.ExecuteActivity(ctx, (*activities.SteamActivities).GetAchievementsForGames, playerID, recentlyPlayedGames).Get(ctx, &gamesAchievements)
	if err != nil {
		return err
	}

	// Comparar achievements
	var newAchievements map[int]*steamclient.GameStats
	err = workflow.ExecuteActivity(ctx, (*activities.SteamActivities).FilterNewAchievements, playerID, gamesAchievements).Get(ctx, &newAchievements)
	if err != nil {
		return err
	}

	// Enviar novas conquistas

	return nil
}
