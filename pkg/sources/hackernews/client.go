package hackernews

import (
	"github.com/notarock/technews-bot/pkg/articles"
	"github.com/peterhellberg/hn"
)

const (
	SOURCE_NAME     = "HACKERNEWS"
	ARTICLES_AMOUNT = 25
)

func FetchLatestTopStories() []articles.Article {
	var articleList []articles.Article

	hn := hn.DefaultClient

	ids, err := hn.TopStories()
	if err != nil {
		panic(err)
	}

	for _, id := range ids[:ARTICLES_AMOUNT] {
		item, err := hn.Item(id)
		if err != nil {
			panic(err)
		}

		article := articles.Article{
			ID:     articles.LinkToID(item.URL),
			Title:  item.Title,
			Link:   item.URL,
			Tags:   []string{},
			Author: item.By,
			Source: SOURCE_NAME,
		}

		articleList = append(articleList, article)
	}

	return articleList
}
