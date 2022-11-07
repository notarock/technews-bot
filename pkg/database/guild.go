package database

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const GUILD_COLLECTION = "discord_guilds"

type Guild struct {
	ID        string   `json:"id,omitempty" bson:"_id,omitempty"`
	GuildID   int      `json:"guildId,omitempty" bson:"guildId,omitempty"`
	Name      string   `json:"name" bson:"name,omitempty"`
	Settings  Settings `json:"settings" bson:"settings,omitempty"`
	ChangedAt int64    `json:"changed_at" bson:"changed_at"`
}

type Settings struct {
	ChannelID int      `json:"channelId" bson:"channelId,omitempty"`
	Interests []string `json:"interests" bson:"interests,omitempty"`
}

func NewGuild(name string, guildID int) *Guild {
	return &Guild{
		GuildID:  guildID,
		Name:     name,
		Settings: Settings{},
	}
}

func InsertGuild(g *Guild) error {
	manyContacts := []interface{}{
		g,
	}

	insertResult, err := collections.Guild.InsertMany(context.TODO(), manyContacts)

	if err != nil {
		log.Panic(err)
	}
	contactIDs := insertResult.InsertedIDs
	var contactIDs_ []primitive.ObjectID

	for _, id := range contactIDs {
		contactIDs_ = append(contactIDs_, id.(primitive.ObjectID))
	}

	fmt.Printf("Inserted %v %T\n", contactIDs_, contactIDs_)

	return err
}
