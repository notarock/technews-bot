package database

import (
	"context"
	"fmt"
	"time"

	"github.com/notarock/technews-bot/pkg/telemetry"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
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

func FindGuildByGuildID(ctx context.Context, guildID string) (Guild, error) {
	ctx, span := telemetry.Tracer.Start(ctx, "database.FindGuildByGuildID")
	defer span.End()

	var guild Guild
	filter := bson.M{"guildId": guildID}
	err := collections.Guild.FindOne(ctx, filter).Decode(&guild)
	if err != nil && err != mongo.ErrNoDocuments {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to find guild")
		return guild, err
	}

	return guild, nil
}

func InsertGuild(ctx context.Context, g Guild) (Guild, error) {
	ctx, span := telemetry.Tracer.Start(ctx, "database.InsertGuild")
	defer span.End()

	guild, err := collections.Guild.InsertOne(ctx, g)
	g.ID = guild.InsertedID.(string)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to insert guild")
	}
	return g, err
}

func GetAllGuilds(ctx context.Context) (guilds []Guild, err error) {
	ctx, span := telemetry.Tracer.Start(ctx, "database.GetAllGuilds")
	defer span.End()

	cur, err := collections.Guild.Find(ctx, bson.D{}, options.Find())
	if err != nil && err != mongo.ErrNoDocuments {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to get all guilds")
		return []Guild{}, err
	}

	for cur.Next(ctx) {
		var g Guild
		err := cur.Decode(&g)
		if err != nil {
			span.RecordError(err)
			return guilds, err
		}

		guilds = append(guilds, g)
	}
	span.SetAttributes(attribute.Int("guilds.count", len(guilds)))
	return guilds, nil
}

func (g Guild) Save(ctx context.Context) error {
	ctx, span := telemetry.Tracer.Start(ctx, "database.Guild.Save")
	defer span.End()

	objectID, err := primitive.ObjectIDFromHex(g.ID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "invalid object ID")
		return err
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": bson.M{
		"guildId":         g.GuildID,
		"name":            g.Name,
		"Settings":        g.Settings,
		"channelSubjects": g.ChannelSubjects,
		"changed_at":      time.Now().Unix(),
	}}
	r, err := collections.Guild.UpdateOne(ctx, filter, update)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to save guild")
	}
	fmt.Printf("%+v\n", g)
	fmt.Printf("%+v\n", r)
	return err
}

func (g *Guild) AddChannelSubject(channelID string, subject string) {
	for i := 0; i < len(g.ChannelSubjects); i++ {
		if g.ChannelSubjects[i].ChannelID == channelID {

			for _, channelSubject := range g.ChannelSubjects[i].Subjects {
				if channelSubject == subject {
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
