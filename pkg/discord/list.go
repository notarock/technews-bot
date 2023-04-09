package discord

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/notarock/technews-bot/pkg/database"
)

func listSubjects(s *discordgo.Session, m *discordgo.MessageCreate) discordgo.MessageEmbed {
	guild, err := database.FindGuildByGuildID(m.GuildID)

	if err != nil {
		return discordgo.MessageEmbed{
			Title:       "Error",
			Description: "Guild not found",
		}
	}

	response := discordgo.MessageEmbed{
		Title:       "Subjects",
		Description: "List of subjects",
	}

	for _, channelSubjects := range guild.ChannelSubjects {
		channel, err := s.Channel(channelSubjects.ChannelID)
		if err != nil {
			continue
		}

		var formated []string

		for _, s := range channelSubjects.Subjects {
			formated = append(formated, fmt.Sprintf("`%s`", s))
		}

		response.Fields = append(response.Fields, &discordgo.MessageEmbedField{
			Name:   channel.Name,
			Value:  strings.Join(formated, "  "),
			Inline: false,
		})
	}

	return response

}
