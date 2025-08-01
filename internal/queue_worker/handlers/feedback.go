package handlers

import (
	"encoding/json"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/rabbitmq/amqp091-go"
	"github.com/v-venes/friends-achievements-bot/pkg/broker"
)

type NewFeedbackHandlerParams struct {
	DiscordSession *discordgo.Session
}

var FeedbackColors map[broker.FeedbackMessageTypeEnum]int = map[broker.FeedbackMessageTypeEnum]int{
	broker.ErrorMessage: 10820909,
	broker.SuccessMessage: 2531945,
}


func NewFeedbackHandler(params NewFeedbackHandlerParams) func(d amqp091.Delivery) {
	return func(d amqp091.Delivery) {
		var feedbackMessage broker.SendFeedbackMessage

		err := json.Unmarshal(d.Body, &feedbackMessage)
		if err != nil {
			log.Printf("Erro ao fazer Unmarshal da mensage: %s", err.Error())
			return
		}

		currentTime := time.Now().Format("02/01/2006 15:04")

		messageEmbed := &discordgo.MessageEmbed{
			Title:       "Novo SteamID",
			Description: feedbackMessage.Content,
			Color: FeedbackColors[feedbackMessage.Type],
			Author: &discordgo.MessageEmbedAuthor{
				Name: "Feedback processamento",
			},
			Footer: &discordgo.MessageEmbedFooter{
				Text: "Tentativa em " + currentTime,
			},
			Fields: []*discordgo.MessageEmbedField{{Name: "Solicitação feita por", Value: feedbackMessage.Username}},
		}

		_, err = params.DiscordSession.ChannelMessageSendEmbed(feedbackMessage.ChannelID, messageEmbed)

		if err != nil {
			log.Printf("Erro ao enviar mensagem [%s]: %s", feedbackMessage.ChannelID, err.Error())
			return
		}

	}
}
