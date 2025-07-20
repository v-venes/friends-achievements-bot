package discordbot

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/v-venes/friends-achievements-bot/internal/discord-bot/commands"
	"github.com/v-venes/friends-achievements-bot/pkg/broker"
)

type Bot struct {
	discordSession *discordgo.Session
	discordGuidId  string
	broker         *broker.Broker
}

type NewBotParams struct {
	DiscordToken   string
	DisocrdGuildID string
	Broker         *broker.Broker
}

func NewBot(params NewBotParams) (*Bot, error) {
	session, err := discordgo.New("Bot " + params.DiscordToken)
	if err != nil {
		return nil, err
	}

	return &Bot{
		discordSession: session,
		discordGuidId:  params.DisocrdGuildID,
		broker:         params.Broker,
	}, nil
}

func (b *Bot) registerHandlers(handlers map[string]commands.CommandHandler) {
	b.discordSession.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Bot started")
	})
	b.discordSession.AddHandler(commands.CreateRouter(handlers))
}

func (b *Bot) registerCommands(slashCommands []commands.SlashCommand) map[string]commands.CommandHandler {
	handlers := make(map[string]commands.CommandHandler)

	for _, cmd := range slashCommands {
		created, err := b.discordSession.ApplicationCommandCreate(b.discordSession.State.User.ID, b.discordGuidId, cmd.Command)
		if err != nil {
			log.Printf("Cannot register '%v' command: %v", cmd.Command.Name, err)
			continue
		}
		handlers[created.Name] = cmd.Handler
	}

	return handlers
}

func (b *Bot) unregisterCommands(commands []commands.SlashCommand) {
	appID := b.discordSession.State.User.ID

	for _, cmd := range commands {
		appCommands, err := b.discordSession.ApplicationCommands(appID, b.discordGuidId)
		if err != nil {
			log.Printf("Erro ao listar comandos: %v", err)
			continue
		}

		for _, registered := range appCommands {
			if registered.Name == cmd.Command.Name {
				err := b.discordSession.ApplicationCommandDelete(appID, b.discordGuidId, registered.ID)
				if err != nil {
					log.Printf("Erro ao deletar comando %s: %v", registered.Name, err)
				} else {
					log.Printf("Comando /%s desregistrado com sucesso", registered.Name)
				}
				break
			}
		}
	}
}

func (b *Bot) Run() error {
	err := b.discordSession.Open()
	if err != nil {
		fmt.Printf("Cannot open the session: %v", err)
		return err
	}

	log.Println("Adding commands...")
	commands := commands.GetSlashCommands(commands.SlashCommandRouterParams{Broker: b.broker})
	handlers := b.registerCommands(commands)
	b.registerHandlers(handlers)

	defer b.discordSession.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	b.unregisterCommands(commands)
	log.Println("Gracefully shutting down")

	return nil
}
