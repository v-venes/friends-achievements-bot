package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var AddAccountCommand *discordgo.ApplicationCommand = &discordgo.ApplicationCommand{
	Name:        "add-account",
	Description: "Add steam account for a user",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Name:        "steamid",
			Description: "Steam account ID",
			Type:        discordgo.ApplicationCommandOptionString,
			Required:    true,
		},
	},
}

func AddAccountHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ApplicationCommandData()
	fmt.Println("Sending response")
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf(
				"Your SteamID is %q",
				data.Options[0].StringValue(),
			),
		},
	})
	if err != nil {
		panic(err)
	}
}
