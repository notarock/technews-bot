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
			Name:   "`list`",
			Value:  "Print all bound channels and their related subjects.",
			Inline: true,
		},
		{
			Name:   "`add`",
			Value:  "Add a technews subject to a bound channel.",
			Inline: true,
		},
		{
			Name:   "`remove`",
			Value:  "Remove a technews subject from a bound channel.",
			Inline: true,
		},
	},
}
