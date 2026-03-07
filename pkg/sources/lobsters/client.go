package lobsters

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/notarock/technews-bot/pkg/articles"
	"github.com/notarock/technews-bot/pkg/telemetry"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

const (
	LOBSTER_URL    = "https://lobste.rs"
	TAGS_SEPARATOR = ","
	SOURCE_NAME    = "LOBSTERS"
)

func FetchLatestArticles(ctx context.Context) []articles.Article {
	var lobsterArticles []articles.Article

	ctx, span := telemetry.Tracer.Start(ctx, "lobsters.FetchLatestArticles")
	defer span.End()

	res, err := http.Get(LOBSTER_URL)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to fetch articles")
		log.Fatalln(err)
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to parse HTML")
		log.Fatal(err)
	}

	span.SetAttributes(attribute.Int("status.code", res.StatusCode))

	doc.Find(".details").Each(func(i int, s *goquery.Selection) {
		ul := s.Find(".u-url")
		title := ul.Text()
		link, _ := ul.Attr("href")
		author := s.Find(".u-author").Text()

		var taglist []string
		if tags, ok := s.Find(".tag").Attr("title"); ok {
			taglist = strings.Split(tags, TAGS_SEPARATOR)
		}

		commentsLink := ""
		s.Find(".comments_label a").EachWithBreak(func(i int, a *goquery.Selection) bool {
			href, exists := a.Attr("href")
			if exists {
				commentsLink = href
				return false
			}
			return true
		})
		if commentsLink != "" && !strings.HasPrefix(commentsLink, "http") {
			commentsLink = LOBSTER_URL + commentsLink
		}

		article := articles.Article{
			ID:         articles.LinkToID(link),
			Title:      title,
			Link:       link,
			Tags:       taglist,
			Author:     author,
			Source:     SOURCE_NAME,
			ThreadLink: commentsLink,
		}

		lobsterArticles = append(lobsterArticles, article)
	})

	span.SetAttributes(attribute.Int("articles.fetched", len(lobsterArticles)))
	return lobsterArticles
}
