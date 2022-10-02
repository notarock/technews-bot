package lobsters

import (
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	LOBSTER_URL    = "https://lobste.rs"
	TAGS_SEPARATOR = ","
)

type LobsterArticle struct {
	Title string
	Link  string
	Tags  []string
}

func FetchLatest() []LobsterArticle {
	var articles []LobsterArticle

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
		// For each item found, get the title
		ul := s.Find(".u-url")
		title := ul.Text()
		link, _ := ul.Attr("href")

		var taglist []string

		if tags, ok := ul.Find(".tag").Attr("href"); !ok {
			taglist = strings.Split(tags, TAGS_SEPARATOR)
		}

		article := LobsterArticle{
			Title: title,
			Link:  link,
			Tags:  taglist,
		}

		articles = append(articles, article)
	})

	return articles
}
