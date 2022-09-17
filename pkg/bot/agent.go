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
}

type Bot struct {
	discordClient *discord.DiscordClient
	store         map[int]discord.Article
}

func Init(config BotConfig) (Bot, error) {
	return Bot{
		discordClient: &config.DiscordClient,
		store:         map[int]discord.Article{},
	}, nil
}

func (b Bot) Serve() {
	ticker := time.NewTicker(5 * time.Minute)
	done := make(chan bool)

	log.Println("Bot is now running in the background.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

	b.SendArticles()

	go b.discordClient.Wait()

	go func() {
		for {
			select {
			case <-done:
				return
			case t := <-ticker.C:
				b.SendArticles()
				fmt.Println("Tick at", t)
			}
		}
	}()

	fmt.Println("Waiting. Ctrl-c to quit")
	<-sc
	ticker.Stop()
	done <- true
	fmt.Println("Ticker stopped")
}

func GetIssues() []*hn.Item {
	i := hackernews.FetchLatest(20)
	var items []*hn.Item
	keywords := []string{"sre", "linux", "breach", "privacy", "speed", "programming"}
	for _, item := range i {
		for _, keyword := range keywords {
			if strings.Contains(strings.ToLower(item.Title), keyword) {
				items = append(items, item)
				break
			}
		}
	}

	return items
}

func (b Bot) SendArticles() {
	items := GetIssues()
	for _, item := range items {
		fromStore, ok := b.store[item.ID]

		if !ok {
			a := discord.Article{
				Title:   item.Title,
				Link:    item.URL,
				Summary: item.Text,
				Author:  item.By,
			}
			b.discordClient.SendArticle(a)

			b.store[item.ID] = a
		} else {
			log.Println("Article already send", fromStore)
		}
	}
}
