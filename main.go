package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/joho/godotenv/autoload"

	"github.com/notarock/technews-bot/pkg/discord"
	"github.com/notarock/technews-bot/pkg/sources/hackernews"
	"github.com/peterhellberg/hn"
)

func main() {
	i := hackernews.FetchLatest(10)
	keywords := []string{"game"}
	for _, item := range i {
		filder(item, keywords)
	}

	discord, err := discord.Init(discord.DiscordConfig{
		Token: os.Getenv("DISCORD_TOKEN"),
	})

	if err != nil {
		log.Fatal(err)
	}

	discord.Do()

}

func filder(item *hn.Item, keywords []string) {
	for _, keyword := range keywords {
		if strings.Contains(strings.ToLower(item.Title), keyword) {
			fmt.Println(item)
		}
	}
}
