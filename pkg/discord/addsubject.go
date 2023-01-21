package discord

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/notarock/technews-bot/pkg/database"
)

func addSubjectToChannel(s *discordgo.Session, m *discordgo.MessageCreate) discordgo.MessageEmbed {
	guild, err := s.Guild(m.GuildID)
	if err != nil {
		return ErrorEmbed
	}

	guilds, err := database.GetAllGuilds()
	if err != nil {
		return ErrorEmbed
	}

	for _, g := range guilds {
		if g.GuildID == guild.ID && g.Settings.ChannelID == m.ChannelID {
			contentWithoutCmd := strings.ReplaceAll(m.Content, "!technews addsubject ", "")
			_, err := database.AddSubjectToGuild(g, contentWithoutCmd)
			if err != nil {
				return ErrorEmbed
			}

			return discordgo.MessageEmbed{
				Title: "Subject added to channel successfully!",
			}
		}
	}

	return discordgo.MessageEmbed{
		Title: "Cannot bind the subject. Use the bind command first",
	}
}
