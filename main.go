package main

import (
	"fmt"
	"strings"

	"github.com/notarock/technews-bot/pkg/sources/hackernews"
	"github.com/peterhellberg/hn"
)

func main() {
	i := hackernews.FetchLatest(10)
	keywords := []string{"game"}
	for _, item := range i {
		filder(item, keywords)
	}
}

func filder(item *hn.Item, keywords []string) {
	for _, keyword := range keywords {
		if strings.Contains(strings.ToLower(item.Title), keyword) {
			fmt.Println(item.Title)
		}
	}

}
