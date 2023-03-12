package bot

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
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
}

func Init(config BotConfig) (Bot, error) {
	return Bot{
		discordClient: &config.DiscordClient,
	}, nil
}

func (b Bot) Serve() {
	var ticker *time.Ticker
	if os.Getenv("DEBUG") == "true" {
		ticker = time.NewTicker(5 * time.Second) // Uncomment for local testing
	} else {
		ticker = time.NewTicker(5 * time.Minute)
	}

	done := make(chan bool)

	var m sync.Mutex

	log.Println("Bot is now running in the background.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

	go b.discordClient.Wait()

	go func() {
		for {
			select {
			case <-done:
				return
			case t := <-ticker.C:
				fmt.Println("Tick at", t)
				go func() {
					if m.TryLock() {
						fmt.Println("Starting to send articles")
						defer m.Unlock()
						b.filterAndSendArticles()
					} else {
						fmt.Println("Job already running. Skipped")
					}
				}()
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

	for _, article := range aggregation {
		dbArticle, err := database.FindArticleByLink(article.Link)
		if err != nil {
			log.Printf("failed to fetch article in db: %v\n", err)
			continue
		}

		if dbArticle.ID != "" { // Article already sent
			if os.Getenv("DEBUG") == "true" {
				log.Printf("Article already sent: %s\n", article.Link)
			}
			continue
		}

		for _, guild := range guilds {

			for _, chanSubject := range guild.ChannelSubjects {
				sent := false
				for _, subject := range chanSubject.Subjects {
					if article.RelatesTo(subject) && !sent {
						b.discordClient.SendArticle(article, chanSubject.ChannelID)
						time.Sleep(250 * time.Millisecond)
						sent = true
					}
				}
			}
		}

		// Save article in db
		err = database.InsertArticle(context.TODO(), database.Article{
			ArticleID: article.ID,
			Link:      article.Link,
			ChangedAt: 0,
		})
	}
}
