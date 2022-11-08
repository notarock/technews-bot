package discord

import "github.com/bwmarrin/discordgo"

var HelpEmbed = discordgo.MessageEmbed{
	URL:   "https://github.com/notarock/technews-bot",
	Type:  "link",
	Title: "Technews-bot Help",
	Fields: []*discordgo.MessageEmbedField{
		&discordgo.MessageEmbedField{
			Name:   "`help`",
			Value:  "Print this help man-page.",
			Inline: true,
		},
		&discordgo.MessageEmbedField{
			Name:   "`bind`",
			Value:  "Bind technews to the current channel. It will then send it's news articles to that channel.",
			Inline: true,
		},
	},
}
