package discord

import "github.com/bwmarrin/discordgo"

var HelpEmbed = discordgo.MessageEmbed{
	URL:   "https://github.com/notarock/technews-bot",
	Type:  "link",
	Title: "Technews-bot Help",
	Fields: []*discordgo.MessageEmbedField{
		{
			Name:   "`help`",
			Value:  "Print this help man-page.",
			Inline: true,
		},
		{
			Name:   "`bind`",
			Value:  "Bind technews to the current channel. It will then send it's news articles to that channel.",
			Inline: true,
		},
		{
			Name:   "`addsubject`",
			Value:  "Add a technews subject to a bound channel.",
			Inline: true,
		},
	},
}
