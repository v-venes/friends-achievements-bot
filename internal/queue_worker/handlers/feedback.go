package handlers

import (
	"encoding/json"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/rabbitmq/amqp091-go"
	"github.com/v-venes/friends-achievements-bot/pkg/broker"
)

type NewFeedbackHandlerParams struct {
	DiscordSession *discordgo.Session
}

func NewFeedbackHandler(params NewFeedbackHandlerParams) func(d amqp091.Delivery) {
	return func(d amqp091.Delivery) {
		var feedbackMessage broker.SendFeedbackMessage

		err := json.Unmarshal(d.Body, &feedbackMessage)
		if err != nil {
			log.Printf("Erro ao fazer Unmarshal da mensage: %s", err.Error())
			return
		}

		_, err = params.DiscordSession.ChannelMessageSend(feedbackMessage.ChannelID, feedbackMessage.Content)
		if err != nil {
			log.Printf("Erro ao enviar mensagem [%s]: %s", feedbackMessage.ChannelID, err.Error())
			return
		}

	}
}
