package main

import (
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"

	"github.com/notarock/technews-bot/pkg/bot"
	"github.com/notarock/technews-bot/pkg/discord"
)

func main() {
	discord, err := discord.Init(discord.DiscordConfig{
		Token: os.Getenv("DISCORD_TOKEN"),
	})

	if err != nil {
		log.Fatal(err)
	}

	b, err := bot.Init(bot.BotConfig{DiscordClient: discord})

	if err != nil {
		log.Fatal(err)
	}

	b.Server()
}
