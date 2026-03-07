package hackernews

import (
	"context"
	"fmt"

	"github.com/notarock/technews-bot/pkg/articles"
	"github.com/notarock/technews-bot/pkg/telemetry"
	"github.com/peterhellberg/hn"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

const (
	SOURCE_NAME     = "HACKERNEWS"
	ARTICLES_AMOUNT = 25
)

func FetchLatestTopStories(ctx context.Context) []articles.Article {
	var articleList []articles.Article

	ctx, span := telemetry.Tracer.Start(ctx, "hackernews.FetchLatestTopStories")
	defer span.End()

	client := hn.DefaultClient

	ids, err := client.TopStories()
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to fetch top stories")
		panic(err)
	}

	span.SetAttributes(attribute.Int("stories.count", len(ids)))

	for _, id := range ids[:ARTICLES_AMOUNT] {
		item, err := client.Item(id)
		if err != nil {
			span.RecordError(err)
			panic(err)
		}

		if item.URL == "" {
			continue
		}

		article := articles.Article{
			ID:         articles.LinkToID(item.URL),
			Title:      item.Title,
			Link:       item.URL,
			Tags:       []string{},
			Author:     item.By,
			Source:     SOURCE_NAME,
			ThreadLink: fmt.Sprintf("https://news.ycombinator.com/item?id=%d", item.ID),
		}

		articleList = append(articleList, article)
	}

	span.SetAttributes(attribute.Int("articles.fetched", len(articleList)))
	return articleList
}
