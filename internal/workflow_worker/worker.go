package workflowworker

import (
	"log"

	"github.com/v-venes/friends-achievements-bot/internal/workflow_worker/activities"
	"github.com/v-venes/friends-achievements-bot/internal/workflow_worker/workflows"
	steamclient "github.com/v-venes/friends-achievements-bot/pkg/steam_client"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

type WorkflowWorker struct {
	worker      worker.Worker
	steamClient *steamclient.SteamClient
}

type NewWorkflowWorkerParams struct {
	Client      client.Client
	SteamClient *steamclient.SteamClient
}

func NewWorkflowWorker(params NewWorkflowWorkerParams) *WorkflowWorker {
	temporalWorker := worker.New(params.Client, "steam-achievements", worker.Options{})

	return &WorkflowWorker{
		worker:      temporalWorker,
		steamClient: params.SteamClient,
	}
}

func (w *WorkflowWorker) Run() {
	w.registerWorkflowsAndActivities()

	err := w.worker.Run(worker.InterruptCh())

	if err != nil {
		log.Fatalf("Erro ao iniciar worker: %s", err.Error())
	}
}

func (w *WorkflowWorker) registerWorkflowsAndActivities() {
	worklfowOptions := workflow.RegisterOptions{
		Name: "ExtractPlayerGames",
	}
	w.worker.RegisterWorkflowWithOptions(workflows.ExtractRecentlyPlayedGamesWorkflow, worklfowOptions)

	steamActivities := &activities.SteamActivities{
		Client: w.steamClient,
	}
	w.worker.RegisterActivity(steamActivities)
}
