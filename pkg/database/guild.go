package database

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const GUILD_COLLECTION = "discord_guilds"

type Guild struct {
	ID        string        `json:"id,omitempty" bson:"_id,omitempty"`
	GuildID   string        `json:"guildId,omitempty" bson:"guildId,omitempty"`
	Name      string        `json:"name" bson:"name,omitempty"`
	Settings  GuildSettings `json:"settings" bson:"settings,omitempty"`
	ChangedAt int64         `json:"changed_at" bson:"changed_at"`
}

type GuildSettings struct {
	ChannelID string   `json:"channelId" bson:"channelId,omitempty"`
	Subjects  []string `json:"subjects" bson:"subjects,omitempty"`
}

func NewGuild(name, guildID string, s GuildSettings) *Guild {
	return &Guild{
		GuildID:  guildID,
		Name:     name,
		Settings: s,
	}
}

func InsertGuild(g *Guild) (*Guild, error) {
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

	return g, err
}

func GetAllGuilds() (guilds []Guild, err error) {
	cur, err := collections.Guild.Find(context.TODO(), bson.D{}, options.Find())
	if err != nil {
		return guilds, err
	}

	for cur.Next(context.TODO()) {
		//Create a value into which the single document can be decoded
		var g Guild
		err := cur.Decode(&g)
		if err != nil {
			return guilds, err
		}

		guilds = append(guilds, g)
	}
	return guilds, nil
}

func AddSubjectToGuild(g Guild, subject string) (int64, error) {
	var selectedGuild Guild
	filter := bson.D{primitive.E{Key: "guildId", Value: g.GuildID}}
	err := collections.Guild.FindOne(context.TODO(), filter).Decode(&selectedGuild)
	if err != nil {
		return 0, err
	}

	selectedGuild.Settings.Subjects = append(selectedGuild.Settings.Subjects, subject)
	update := bson.D{primitive.E{Key: "$set", Value: bson.D{primitive.E{Key: "settings", Value: selectedGuild.Settings}}}}
	result, err := collections.Guild.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return 0, err
	}

	return result.MatchedCount, nil // the result.MatchCount should equal to 1
}
