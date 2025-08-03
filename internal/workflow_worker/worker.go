package workflowworker

import (
	"log"

	"github.com/v-venes/friends-achievements-bot/internal/workflow_worker/activities"
	"github.com/v-venes/friends-achievements-bot/internal/workflow_worker/workflows"
	"github.com/v-venes/friends-achievements-bot/pkg/repository"
	steamclient "github.com/v-venes/friends-achievements-bot/pkg/steam_client"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

type WorkflowWorker struct {
	worker           worker.Worker
	steamClient      *steamclient.SteamClient
	playerRepository *repository.PlayerRepository
}

type NewWorkflowWorkerParams struct {
	Client           client.Client
	SteamClient      *steamclient.SteamClient
	PlayerRepository *repository.PlayerRepository
}

func NewWorkflowWorker(params NewWorkflowWorkerParams) *WorkflowWorker {
	temporalWorker := worker.New(params.Client, "steam-achievements", worker.Options{})

	return &WorkflowWorker{
		worker:           temporalWorker,
		steamClient:      params.SteamClient,
		playerRepository: params.PlayerRepository,
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
	w.worker.RegisterWorkflowWithOptions(workflows.ExtractNewPlayerDataWorkflow, workflow.RegisterOptions{
		Name: "ExtractNewPlayerData",
	})
	w.worker.RegisterWorkflowWithOptions(workflows.ExtractGameDataWorkflow, workflow.RegisterOptions{
		Name: "ExtractGameData",
	})

	steamActivities := &activities.SteamActivities{
		Client: w.steamClient,
	}
	playerActitivities := &activities.PlayerActivities{
		PlayerRepository: w.playerRepository,
	}

	w.worker.RegisterActivity(steamActivities)
	w.worker.RegisterActivity(playerActitivities)
}
