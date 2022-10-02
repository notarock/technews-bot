package discord

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/joho/godotenv/autoload"

	"github.com/bwmarrin/discordgo"
	"github.com/notarock/technews-bot/pkg/articles"
)

type DiscordConfig struct {
	Token   string
	Channel string
}

type DiscordClient struct {
	client  *discordgo.Session
	channel string
}

func Init(config DiscordConfig) (DiscordClient, error) {

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		return DiscordClient{}, fmt.Errorf("error creating Discord session: %v", err)
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(healthcheckHandler)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	return DiscordClient{
		client:  dg,
		channel: config.Channel,
	}, nil

	// // Open a websocket connection to Discord and begin listening.
}
func (dc DiscordClient) Wait() error {
	err := dc.client.Open()
	if err != nil {
		return fmt.Errorf("error opening connection: %v", err)
	}

	// Wait here until CTRL-C or other term signal is received.
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dc.client.Close()
	return nil
}

func healthcheckHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	// If the message is "ping" reply with "Pong!"
	if m.Content == "!ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}
}

func (dc DiscordClient) SendArticle(a articles.Article) {
	log.Println(fmt.Sprintf("Attempting to send article named %+v", a))
	embed := a.ToDiscordEmbed()
	m, _ := dc.client.ChannelMessageSendEmbed(dc.channel, embed)
	log.Println(m)
}
