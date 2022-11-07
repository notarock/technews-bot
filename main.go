package main

import (
	"context"
	"os"

	_ "github.com/joho/godotenv/autoload"

	"github.com/notarock/technews-bot/pkg/bot"
	"github.com/notarock/technews-bot/pkg/database"
	"github.com/notarock/technews-bot/pkg/discord"
	log "github.com/sirupsen/logrus"
)

func main() {
	mongodbConfig := database.MongodbConfig{
		Uri:    os.Getenv("MONGODB_URI"),
		DbName: os.Getenv("MONGODB_DBNAME"),
	}

	err := database.Connect(mongodbConfig)
	if err != nil {
		log.Fatal(err)
	}

	err = database.Client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	discordConfig := discord.DiscordConfig{
		Token:   os.Getenv("DISCORD_TOKEN"),
		Channel: os.Getenv("DISCORD_CHANNEL"),
	}

	if discordConfig.Channel == "" || discordConfig.Token == "" {
		log.Fatalln("Cant start technews bot: some env variables are missing.")
	}

	discord, err := discord.Init(discordConfig)

	if err != nil {
		log.Fatal(err)
	}

	b, err := bot.Init(bot.BotConfig{DiscordClient: discord})

	if err != nil {
		log.Fatal(err)
	}

	b.Serve()
}
