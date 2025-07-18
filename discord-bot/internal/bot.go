package internal

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/v-venes/friends-achievements-bot/discord-bot/internal/commands"
)

var (
	botCommands = []*discordgo.ApplicationCommand{
		commands.AddAccountCommand,
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"add-account": commands.AddAccountHandler,
	}
)

type Bot struct {
	botSession *discordgo.Session
}

func NewBot(token string) (*Bot, error) {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}

	return &Bot{
		botSession: session,
	}, nil
}

func configureCommands(session *discordgo.Session) {
	session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Bot started")
	})
	session.AddHandler(func(s *discordgo.Session, r *discordgo.MessageInteractionMetadata) {
	})
	session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		handler, commandFound := commandHandlers[i.ApplicationCommandData().Name]

		fmt.Printf("commandFound: %t", commandFound)
		if commandFound {
			handler(s, i)
		}
	})
}

func (b *Bot) StartBot() error {
	configureCommands(b.botSession)
	err := b.botSession.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	defer b.botSession.Close()

	log.Println("Adding commands...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(botCommands))

	for i, v := range botCommands {
		cmd, err := b.botSession.ApplicationCommandCreate(b.botSession.State.User.ID, "", v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("Gracefully shutting down")

	for _, cmd := range registeredCommands {
		err := b.botSession.ApplicationCommandDelete(b.botSession.State.User.ID, "", cmd.ID)
		if err != nil {
			log.Fatalf("Cannot delete %q command: %v", cmd.Name, err)
		}
	}

	return nil
}
