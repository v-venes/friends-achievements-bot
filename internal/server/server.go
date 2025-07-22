package server

import (
	"encoding/json"
	"fmt"

	"github.com/v-venes/friends-achievements-bot/internal/discord-bot/shared/models"
	"github.com/v-venes/friends-achievements-bot/pkg/broker"
	"github.com/v-venes/friends-achievements-bot/pkg/service"
)

type Server struct {
	broker      *broker.Broker
	steamClient *service.SteamClient
}

type NewServerParams struct {
	Broker      *broker.Broker
	SteamClient *service.SteamClient
}

func NewServer(params NewServerParams) *Server {
	return &Server{
		broker:      params.Broker,
		steamClient: params.SteamClient,
	}
}

func (s *Server) Run() error {
	msgs, err := s.broker.ReceiveMessages(broker.BrokerQueues[broker.NewSteamId])
	if err != nil {
		fmt.Printf("Cannot open the session: %v", err)
		return err
	}

	var forever chan struct{}

	go func() {
		for message := range msgs {
			var addAccountMessage models.AddAccountMessage
			err := json.Unmarshal(message.Body, &addAccountMessage)
			if err != nil {
				fmt.Printf("Erro ao fazer Unmarshal da mensagem: %s", err.Error())
				continue
			}

			_, err = s.steamClient.GetPlayerSummary(addAccountMessage.SteamID)
			if err != nil {
				fmt.Printf("%s", err.Error())
				continue
			}

			// TODO: Salvar Player no banco
			// TODO: Enviar feedback no discord
		}
	}()

	<-forever
	return nil
}
