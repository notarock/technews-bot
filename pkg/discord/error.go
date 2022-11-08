package discord

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var ErrorEmbed = discordgo.MessageEmbed{
	URL:         "https://github.com/notarock/technews-bot",
	Title:       "Invalid command",
	Description: "Error occured... let us know.",
	Fields:      []*discordgo.MessageEmbedField{},
}

func invalidCommandError(command string) discordgo.MessageEmbed {
	return discordgo.MessageEmbed{
		URL:         "https://github.com/notarock/technews-bot",
		Title:       "Invalid command",
		Description: fmt.Sprintf("I did not recognize the command `%s`... are you sure it's valid?", command),
		Fields:      []*discordgo.MessageEmbedField{},
	}

}
