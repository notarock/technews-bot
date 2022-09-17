package bot

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/notarock/technews-bot/pkg/discord"
	"github.com/notarock/technews-bot/pkg/sources/hackernews"
	"github.com/peterhellberg/hn"
)

type BotConfig struct {
	DiscordClient discord.DiscordClient
	store         map[string]string
}

type Bot struct {
	discordClient *discord.DiscordClient
}

func Init(config BotConfig) (Bot, error) {
	return Bot{
		discordClient: &config.DiscordClient,
	}, nil
}

func (b Bot) Server() {
	ticker := time.NewTicker(500 * time.Millisecond)
	done := make(chan bool)

	log.Println("Bot is now running in the background.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

	go b.discordClient.Wait()

	go func() {
		for {
			select {
			case <-sc:
				return
			case <-done:
				return
			case t := <-ticker.C:
				items := GetIssues()
				for _, item := range items {
					b.discordClient.SendMessage(item.Title)

				}
				fmt.Println("Tick at", t)
			}
		}
	}()

	time.Sleep(1600 * time.Millisecond)
	ticker.Stop()
	done <- true
	fmt.Println("Ticker stopped")

}

func GetIssues() []*hn.Item {
	i := hackernews.FetchLatest(10)
	var items []*hn.Item
	keywords := []string{"game"}
	for _, item := range i {
		for _, keyword := range keywords {
			if strings.Contains(strings.ToLower(item.Title), keyword) {
				items = append(items, item)
			}
		}
	}

	return items
}
