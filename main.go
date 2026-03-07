package main

import (
	"context"
	"os"

	_ "github.com/joho/godotenv/autoload"

	"github.com/notarock/technews-bot/pkg/bot"
	"github.com/notarock/technews-bot/pkg/database"
	"github.com/notarock/technews-bot/pkg/discord"
	"github.com/notarock/technews-bot/pkg/telemetry"
	log "github.com/sirupsen/logrus"
)

func main() {
	ctx := context.Background()

	otlpEndpoint := os.Getenv("OTLP_ENDPOINT")
	if otlpEndpoint == "" {
		otlpEndpoint = "localhost:4317"
	}

	shutdown, err := telemetry.Init(ctx, "technews-bot", otlpEndpoint)
	if err != nil {
		log.Warnf("Failed to initialize telemetry: %v", err)
	} else {
		defer shutdown(ctx)
	}

	mongodbConfig := database.MongodbConfig{
		Uri:    os.Getenv("MONGODB_URI"),
		DbName: os.Getenv("MONGODB_DBNAME"),
	}

	err = database.Connect(mongodbConfig)
	if err != nil {
		log.Fatal(err)
	}

	err = database.Healthcheck()
	if err != nil {
		log.Fatal(err)
	}

	discordConfig := discord.DiscordConfig{
		Token: os.Getenv("DISCORD_TOKEN"),
	}

	if discordConfig.Token == "" {
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
