package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/notarock/technews-bot/pkg/database"
	log "github.com/sirupsen/logrus"
)

func bindToChannel(s *discordgo.Session, m *discordgo.MessageCreate) discordgo.MessageEmbed {
	//TODO: Check if guild exists in database before adding

	guild, err := s.Guild(m.GuildID)
	if err != nil {
		return ErrorEmbed
	}

	dbGuild, err := database.InsertGuild(database.NewGuild(guild.Name, guild.ID, database.GuildSettings{
		ChannelID: m.ChannelID,
	}))
	if err != nil {
		return ErrorEmbed
	}

	log.Printf("%+v\n", dbGuild)

	return discordgo.MessageEmbed{
		Title: "Bound to channel successfully!",
	}
}
