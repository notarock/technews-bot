package bot

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/notarock/technews-bot/pkg/articles"
	"github.com/notarock/technews-bot/pkg/discord"
	"github.com/notarock/technews-bot/pkg/sources/hackernews"
	"github.com/notarock/technews-bot/pkg/sources/lobsters"
)

type BotConfig struct {
	DiscordClient discord.DiscordClient
}

type Bot struct {
	discordClient *discord.DiscordClient
	store         map[string]articles.Article
}

func Init(config BotConfig) (Bot, error) {
	return Bot{
		discordClient: &config.DiscordClient,
		store:         map[string]articles.Article{},
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

func GetFilteredIssues() []articles.Article {
	aggregation := aggregateArticles()

	var filteredArticles []articles.Article
	subjects := []string{"sre", "linux", "breach", "privacy", "speed", "programming", "golang", "development"}

	for _, article := range aggregation {
		for _, subject := range subjects {
			if article.RelatesTo(subject) {
				filteredArticles = append(filteredArticles, article)
				break
			}
		}
	}

	return filteredArticles
}

func aggregateArticles() []articles.Article {
	var aggregation []articles.Article

	hnArticles := hackernews.FetchLatestTopStories()
	aggregation = append(aggregation, hnArticles...)

	lbArticles := lobsters.FetchLatestArticles()
	aggregation = append(aggregation, lbArticles...)

	return aggregation
}

func (b Bot) SendArticles() {
	filteredArticles := GetFilteredIssues()
	for _, article := range filteredArticles {
		fromStore, ok := b.store[article.ID]

		if !ok {
			b.discordClient.SendArticle(article)
			b.store[article.ID] = article
		} else {
			log.Println("Article already send", fromStore)
		}
	}
}
