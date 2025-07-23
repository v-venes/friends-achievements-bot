package handlers

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/rabbitmq/amqp091-go"
	"github.com/v-venes/friends-achievements-bot/pkg/broker"
	"github.com/v-venes/friends-achievements-bot/pkg/models"
	"github.com/v-venes/friends-achievements-bot/pkg/service"
)

type NewSteamIDHandlerParams struct {
	// Mongo *MongoClient
	Broker      *broker.Broker
	SteamClient *service.SteamClient
}

func NewSteamIDHandler(params NewSteamIDHandlerParams) func(d amqp091.Delivery) {
	return func(d amqp091.Delivery) {
		var addAccountMessage models.AddAccountMessage
		err := json.Unmarshal(d.Body, &addAccountMessage)
		if err != nil {
			fmt.Printf("Erro ao fazer Unmarshal da mensagem: %s", err.Error())
			return
		}

		_, err = params.SteamClient.GetPlayerSummary(addAccountMessage.SteamID)
		if err != nil {
			fmt.Printf("%s", err.Error())
			return
		}

		// TODO: Salvar Player no banco
		payload := models.SendFeedbackMessage{
			Content:    "",
			Username:   "",
			GuildID:    "",
			ChannelID:  "",
			ExecutedAt: time.Now(),
		}

		err = params.Broker.SendMessage(broker.SendMessageParams{
			Queue:   broker.NewSteamId,
			Message: payload,
		})
		if err != nil {
			return
		}
	}
}
