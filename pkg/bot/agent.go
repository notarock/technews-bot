package bot

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/notarock/technews-bot/pkg/articles"
	"github.com/notarock/technews-bot/pkg/database"
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

	b.filterAndSendArticles()

	go b.discordClient.Wait()

	go func() {
		for {
			select {
			case <-done:
				return
			case t := <-ticker.C:
				b.filterAndSendArticles()
				fmt.Println("Tick at", t)
			}
		}
	}()

	fmt.Println("Waiting. Ctrl-c to quit")
	<-sc
	ticker.Stop()
	done <- true
	fmt.Println("Received closing signal: exiting.")
}

func aggregateArticles() []articles.Article {
	var aggregation []articles.Article

	hnArticles := hackernews.FetchLatestTopStories()
	aggregation = append(aggregation, hnArticles...)

	lbArticles := lobsters.FetchLatestArticles()
	aggregation = append(aggregation, lbArticles...)

	return aggregation
}

func (b Bot) filterAndSendArticles() {
	if os.Getenv("DRY_RUN") == "true" {
		log.Println("Dry run, do nothing")
		return
	}

	guilds, err := database.GetAllGuilds()
	if err != nil {
		log.Println(err)
		return
	}

	aggregation := aggregateArticles()

	for _, guild := range guilds {
		filteredArticles := filterArticles(aggregation, guild)
		for _, article := range filteredArticles {
			fromStore, ok := b.store[article.ID]

			if !ok {
				b.discordClient.SendArticle(article, guild.Settings.ChannelID)
				b.store[article.ID] = article
				// Don't wanna spam too much aight
				time.Sleep(3 * time.Second)
			} else {
				log.Println("Article already sent", fromStore)
			}
		}
	}
}

func filterArticles(aggregation []articles.Article, guild database.Guild) []articles.Article {
	var filteredArticles []articles.Article
	for _, article := range aggregation {
		for _, subject := range guild.Settings.Subjects {
			if article.RelatesTo(subject) {
				filteredArticles = append(filteredArticles, article)
				break // Do not send article twice when match many subjects
			}
		}
	}

	return filteredArticles
}
