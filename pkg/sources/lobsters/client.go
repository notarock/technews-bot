package lobsters

import (
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/notarock/technews-bot/pkg/articles"
)

const (
	LOBSTER_URL    = "https://lobste.rs"
	TAGS_SEPARATOR = ","
	SOURCE_NAME    = "LOBSTERS"
)

func FetchLatestArticles() []articles.Article {
	var lobsterArticles []articles.Article

	res, err := http.Get(LOBSTER_URL)
	if err != nil {
		log.Fatalln(err)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find(".details").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the title, link, tags and author.
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
				return false // break after first match
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

	return lobsterArticles
}
