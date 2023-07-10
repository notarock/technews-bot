package discord

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/notarock/technews-bot/pkg/database"
	log "github.com/sirupsen/logrus"
)

func addSubjectToChannel(s *discordgo.Session, m *discordgo.MessageCreate) discordgo.MessageEmbed {
	subjectToAdd := strings.ReplaceAll(m.Content, "!technews add ", "")
	guild, err := database.FindGuildByGuildID(m.GuildID)
	if err != nil {
		log.Errorf("error occured while trying to find guild by guildId: %v:", err)
		return ErrorEmbed
	}

	if guild.ID == "" { // Guild not found
		discordGuild, err := s.Guild(m.GuildID)
		if err != nil {
			log.Errorf("error occured while trying to find discord guild: %v:", err)
			return ErrorEmbed
		}

		guild, err = database.InsertGuild(database.NewGuild(discordGuild.ID, discordGuild.Name))

		if err != nil {
			log.Errorf("error occured while trying to insert new guild in db: %v:", err)
			return ErrorEmbed
		}
	}

	guild.AddChannelSubject(m.ChannelID, subjectToAdd)
	err = guild.Save()
	if err != nil {
		log.Errorf("error occured while trying to save subject in db: %v:", err)
		return ErrorEmbed
	}

	return discordgo.MessageEmbed{
		Title: fmt.Sprintf("Subject %s added to channel successfully!", subjectToAdd),
	}

}
