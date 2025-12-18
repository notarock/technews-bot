package discord

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	_ "github.com/joho/godotenv/autoload"

	"github.com/bwmarrin/discordgo"
	"github.com/notarock/technews-bot/pkg/articles"
	"github.com/notarock/technews-bot/pkg/summarize"
)

type DiscordConfig struct {
	Token string
}

type DiscordClient struct {
	client *discordgo.Session
}

func Init(config DiscordConfig) (DiscordClient, error) {

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		return DiscordClient{}, fmt.Errorf("error creating Discord session: %v", err)
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageHandler)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	return DiscordClient{
		client: dg,
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

// messageHandler handles an incomming message, checks for command action and execute them
// if they are prefixed with "!technews".
func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Handles basic !technews commands
	if strings.HasPrefix(m.Content, "!technews") {
		command := parseCommandMessage(m.Content)
		fmt.Printf("received command %+v\n", command)

		response := command.Execute(s, m)

		_, err := s.ChannelMessageSendEmbed(m.ChannelID, &response)
		if err != nil {
			log.Println(err)
		}
	}

	// IF reply to the bot with "tldr", summarize the article in the referenced message
	if m.MessageReference != nil && m.ReferencedMessage.Author.ID == s.State.User.ID && strings.ToLower(m.Content) == "tldr" {
		SummarizeRepliedMessage(s, m)
	}
}

func SummarizeRepliedMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	articleURL := m.ReferencedMessage.Embeds[0].URL
	if articleURL == "" {
		return
	}

	// SIgnal that bot is working on it
	s.ChannelTyping(m.ChannelID)
	// No need to defer this, as the typing indicator will stop when a message is sent

	gemini, err := summarize.GetClient()
	if err != nil {
		log.Println(err)
		return
	}

	summary, err := gemini.SummarizeWebpage(articleURL)
	if err != nil {
		log.Println(err)
		return
	}

	_, err = s.ChannelMessageSend(m.ChannelID, summary)
	if err != nil {
		log.Println(err)
	}
}

func (dc DiscordClient) SendArticle(a articles.Article, channel string) {
	log.Printf("Attempting to send article named %+v\n", a)
	embed := a.ToDiscordEmbed()
	m, _ := dc.client.ChannelMessageSendEmbed(channel, embed)
	log.Println(m)
}
