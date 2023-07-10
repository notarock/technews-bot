package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const GUILD_COLLECTION = "discord_guilds"

type Guild struct {
	ID              string            `json:"id,omitempty" bson:"_id,omitempty"`
	GuildID         string            `json:"guildId,omitempty" bson:"guildId,omitempty"`
	Name            string            `json:"name" bson:"name,omitempty"`
	Settings        Settings          `json:"Settings" bson:"Settings,omitempty"`
	ChannelSubjects []ChannelSubjects `json:"channelSubjects" bson:"channelSubjects,omitempty"`
	ChangedAt       int64             `json:"changed_at" bson:"changed_at"`
}

type ChannelSubjects struct {
	ChannelID string   `json:"channelId" bson:"channelId,omitempty"`
	Subjects  []string `json:"subjects" bson:"subjects,omitempty"`
}

type Settings struct {
	Active bool `json:"active" bson:"active"`
}

func NewGuild(guildID, name string) Guild {
	return Guild{
		GuildID: guildID,
		Name:    name,
		Settings: Settings{
			Active: true,
		},
		ChangedAt: time.Now().Unix(),
	}
}

func FindGuildByGuildID(guildID string) (Guild, error) {
	var guild Guild
	filter := bson.M{"guildId": guildID}
	err := collections.Guild.FindOne(context.TODO(), filter).Decode(&guild)
	if err != nil && err != mongo.ErrNoDocuments {
		return guild, err
	}

	return guild, nil
}

func InsertGuild(g Guild) (Guild, error) {
	guild, err := collections.Guild.InsertOne(context.TODO(), g)
	g.ID = guild.InsertedID.(string)
	return g, err
}

func GetAllGuilds() (guilds []Guild, err error) {
	cur, err := collections.Guild.Find(context.TODO(), bson.D{}, options.Find())
	if err != nil && err != mongo.ErrNoDocuments {
		return []Guild{}, err
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

func (g Guild) Save() error {
	objectID, err := primitive.ObjectIDFromHex(g.ID)

	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": bson.M{
		"guildId":         g.GuildID,
		"name":            g.Name,
		"Settings":        g.Settings,
		"channelSubjects": g.ChannelSubjects,
		"changed_at":      time.Now().Unix(),
	}}
	r, err := collections.Guild.UpdateOne(context.TODO(), filter, update)
	fmt.Printf("%+v\n", g)
	fmt.Printf("%+v\n", r)
	return err
}

func (g *Guild) AddChannelSubject(channelID string, subject string) {
	for i := 0; i < len(g.ChannelSubjects); i++ {
		if g.ChannelSubjects[i].ChannelID == channelID {

			for _, channelSubject := range g.ChannelSubjects[i].Subjects {
				if channelSubject == subject { // subject already exists
					return
				}
			}

			g.ChannelSubjects[i].Subjects = append(g.ChannelSubjects[i].Subjects, subject)
			return
		}
	}

	g.ChannelSubjects = append(g.ChannelSubjects, ChannelSubjects{
		ChannelID: channelID,
		Subjects:  []string{subject},
	})
}

func (g *Guild) RemoveChannelSubject(channelID string, subject string) {
	subjectsToKeep := []string{}

	for i := 0; i < len(g.ChannelSubjects); i++ {
		if g.ChannelSubjects[i].ChannelID == channelID {
			subjects := g.ChannelSubjects[i].Subjects

			for _, channelSubject := range subjects {
				if channelSubject != subject {
					subjectsToKeep = append(subjectsToKeep, channelSubject)
				}
			}

			g.ChannelSubjects[i].Subjects = subjectsToKeep
		}
	}
}
