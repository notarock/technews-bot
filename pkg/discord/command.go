package discord

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

type Command struct {
	Action    string
	Arguments []string
}

const (
	COMMAND_TYPE_HELP = "help"
	COMMAND_TYPE_BIND = "bind"
)

func parseCommandMessage(message string) Command {
	parts := strings.Split(message, " ")

	if len(parts) < 2 {
		return Command{
			Action:    "help",
			Arguments: []string{},
		}
	}

	return Command{
		Action:    strings.ToLower(parts[1]),
		Arguments: parts[2:],
	}
}

func (c Command) Execute(s *discordgo.Session, m *discordgo.MessageCreate) discordgo.MessageEmbed {
	switch c.Action {
	case COMMAND_TYPE_HELP:
		return HelpEmbed
	case COMMAND_TYPE_BIND:
		return bindToChannel(s, m)
	default:
		return invalidCommandError(c.Action)
	}
}
