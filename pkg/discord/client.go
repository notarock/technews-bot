package discord

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/joho/godotenv/autoload"

	"github.com/bwmarrin/discordgo"
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

type Article struct {
	Title   string
	Link    string
	Summary string
	Author  string
}

func (dc DiscordClient) SendArticle(a Article) {
	log.Println("Attempting to send article named ", a)

	dc.client.ChannelMessageSendEmbed(dc.channel, &discordgo.MessageEmbed{
		URL:         a.Link,
		Type:        "link",
		Title:       a.Title,
		Description: a.Summary,
		Timestamp:   "",
		Color:       0,
		Author:      &discordgo.MessageEmbedAuthor{URL: fmt.Sprintf("https://news.ycombinator.com/user?id=%s", a.Author), Name: a.Author, IconURL: "https://news.ycombinator.com/y18.gif"},
		Fields:      []*discordgo.MessageEmbedField{},
	})
}
