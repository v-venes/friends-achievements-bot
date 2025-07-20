package commands

import (
	"github.com/bwmarrin/discordgo"
)

type CommandHandler func(s *discordgo.Session, i *discordgo.InteractionCreate)

func CreateRouter(handlers map[string]CommandHandler) func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		name := i.ApplicationCommandData().Name
		if handler, ok := handlers[name]; ok {
			handler(s, i)
		}
	}
}
