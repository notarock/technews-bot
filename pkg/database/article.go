package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const ARTICLE_COLLECTION = "articles_cache"

type Article struct {
	ID        string `json:"id,omitempty" bson:"_id,omitempty"`
	ArticleID string `json:"article_id,omitempty" bson:"article_id,omitempty"`
	Link      string `json:"link,omitempty" bson:"link,omitempty"`
	ChangedAt int64  `json:"changed_at" bson:"changed_at"`
}

func InsertArticle(ctx context.Context, article Article) error {
	_, err := collections.Article.InsertOne(ctx, article)
	return err
}

func FindArticleByLink(link string) (*Article, error) {
	var article Article
	err := collections.Article.FindOne(context.TODO(), bson.M{"link": link}).Decode(&article)

	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}
	return &article, nil
}
