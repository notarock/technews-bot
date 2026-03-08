package database

import (
	"context"

	"github.com/notarock/technews-bot/pkg/telemetry"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.opentelemetry.io/otel/codes"
)

const ARTICLE_COLLECTION = "articles_cache"

type Article struct {
	ID        string `json:"id,omitempty" bson:"_id,omitempty"`
	ArticleID string `json:"article_id,omitempty" bson:"article_id,omitempty"`
	Link      string `json:"link,omitempty" bson:"link,omitempty"`
	ChangedAt int64  `json:"changed_at" bson:"changed_at"`
}

func InsertArticle(ctx context.Context, article Article) error {
	ctx, span := telemetry.Tracer.Start(ctx, "database.InsertArticle")
	defer span.End()

	_, err := collections.Article.InsertOne(ctx, article)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to insert article")
	}
	return err
}

func FindArticleByLink(ctx context.Context, link string) (*Article, error) {
	ctx, span := telemetry.Tracer.Start(ctx, "database.FindArticleByLink")
	defer span.End()

	var article Article
	err := collections.Article.FindOne(ctx, bson.M{"link": link}).Decode(&article)

	if err != nil && err != mongo.ErrNoDocuments {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to find article")
		return nil, err
	}
	return &article, nil
}
