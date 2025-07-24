package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/rabbitmq/amqp091-go"
	"github.com/v-venes/friends-achievements-bot/pkg/broker"
	"github.com/v-venes/friends-achievements-bot/pkg/repository"
	steamclient "github.com/v-venes/friends-achievements-bot/pkg/steam_client"
)

type NewSteamIDHandlerParams struct {
	PlayerRepository *repository.PlayerRepository
	Broker           *broker.Broker
	SteamClient      *steamclient.SteamClient
}

func NewSteamIDHandler(params NewSteamIDHandlerParams) func(d amqp091.Delivery) {
	return func(d amqp091.Delivery) {
		var addAccountMessage broker.AddAccountMessage
		err := json.Unmarshal(d.Body, &addAccountMessage)
		if err != nil {
			log.Printf("Erro ao fazer Unmarshal da mensagem: %s", err.Error())
			return
		}

		var payload broker.SendFeedbackMessage

		steamPlayer, err := params.SteamClient.GetPlayerSummary(addAccountMessage.SteamID)
		if err != nil {
			log.Printf("Erro ao adicionar SteamID [%s] %s\n", addAccountMessage.SteamID, err.Error())
			payload = broker.SendFeedbackMessage{
				Content:    fmt.Sprintf("SteamID [%s] Não foi encontrado!", addAccountMessage.SteamID),
				Type:       broker.ErrorMessage,
				Username:   addAccountMessage.Username,
				GuildID:    addAccountMessage.GuildID,
				ChannelID:  addAccountMessage.ChannelID,
				ExecutedAt: time.Now(),
			}
		} else {
			playerModel := repository.NewPlayerFromSteam(steamPlayer)
			err = params.PlayerRepository.CreatePlayer(*playerModel)
			if err != nil {
				log.Printf("Erro ao criar player %s\n", err.Error())
				return
			}

			payload = broker.SendFeedbackMessage{
				Content:    fmt.Sprintf("SteamID [%s] adicionado com sucesso, começando a extração", addAccountMessage.SteamID),
				Type:       broker.SuccessMessage,
				Username:   addAccountMessage.Username,
				GuildID:    addAccountMessage.GuildID,
				ChannelID:  addAccountMessage.ChannelID,
				ExecutedAt: time.Now(),
			}
		}

		err = params.Broker.SendMessage(broker.SendMessageParams{
			Queue:   broker.FeedbackMessages,
			Message: payload,
		})
		if err != nil {
			log.Printf("%s\n", err.Error())
			return
		}
	}
}
