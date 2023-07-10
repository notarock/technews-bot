package discord

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/notarock/technews-bot/pkg/database"
	log "github.com/sirupsen/logrus"
)

func removeSubjectFromChannel(s *discordgo.Session, m *discordgo.MessageCreate) discordgo.MessageEmbed {
	//TODO: Truncate command from message before reaching this point
	subjectToRemove := strings.ReplaceAll(m.Content, "!technews remove ", "")
	guild, err := database.FindGuildByGuildID(m.GuildID)
	if err != nil {
		log.Errorf("error occured while trying to find guild by guildId: %v:", err)
		return ErrorEmbed
	}

	guild.RemoveChannelSubject(m.ChannelID, subjectToRemove)
	guild.Save()

	return discordgo.MessageEmbed{
		Title: fmt.Sprintf("Subject %s removed to channel successfully!", subjectToRemove),
	}
}
