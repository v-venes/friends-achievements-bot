package commands

import (
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/v-venes/friends-achievements-bot/pkg/broker"
)

type SlashCommand struct {
	Command *discordgo.ApplicationCommand
	Handler CommandHandler
}

type SlashCommandRouterParams struct {
	Broker *broker.Broker
}

type AddAccountPayload struct {
	SteamID    string    `json:"steam_id"`
	Username   string    `json:"username"`
	GuildID    string    `json:"guild_id"`
	ExecutedAt time.Time `json:"executed_at"`
}

func GetSlashCommands(params SlashCommandRouterParams) []SlashCommand {
	return []SlashCommand{
		{
			Command: &discordgo.ApplicationCommand{
				Name:        "adicionar-conta",
				Description: "Adicionar uma conta Steam para consulta",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Name:        "steamid",
						Description: "ID da sua conta Steam",
						Type:        discordgo.ApplicationCommandOptionString,
						Required:    true,
					},
				},
			},
			Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
				data := i.ApplicationCommandData()

				if len(data.Options) == 0 {
					log.Printf("SteamID não encontrado")
					return
				}

				steamID := data.Options[0].StringValue()

				payload := AddAccountPayload{
					SteamID:    steamID,
					Username:   i.Member.User.Username,
					GuildID:    i.GuildID,
					ExecutedAt: time.Now(),
				}

				params.Broker.SendMessage(broker.SendMessageParams{
					Queue:   broker.NewSteamId,
					Message: payload,
				})

				err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: fmt.Sprintf(
							"SteamID [%q] enviado para processamento...",
							steamID,
						),
					},
				})
				if err != nil {
					log.Printf("Erro ao enviar resposta")
					return
				}
			},
		},
	}
}
