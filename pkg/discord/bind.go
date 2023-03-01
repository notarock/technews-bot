package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/notarock/technews-bot/pkg/database"
	log "github.com/sirupsen/logrus"
)

func bindToChannel(s *discordgo.Session, m *discordgo.MessageCreate) discordgo.MessageEmbed {
	guild, err := s.Guild(m.GuildID)
	if err != nil {
		return ErrorEmbed
	}

	dbGuild, err := database.FindGuildByGuildID(guild.ID)

	if dbGuild.ID != "" {
		log.Printf("Guild already exist: %+v\n", dbGuild)
		return discordgo.MessageEmbed{
			Title: "Already bound to channel in this guild!",
		}
	}
	//TODO: Add channel binding to guild and keep a list of channels with their respective subjects

	_, err = database.InsertGuild(database.NewGuild(guild.Name, guild.ID, database.GuildSettings{
		ChannelID: m.ChannelID,
	}))
	if err != nil {
		return ErrorEmbed
	}

	log.Printf("New guild created: %+v\n", dbGuild)

	return discordgo.MessageEmbed{
		Title: "Bound to channel successfully!",
	}

}
