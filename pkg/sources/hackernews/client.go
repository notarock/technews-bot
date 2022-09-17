package hackernews

import (
	"github.com/peterhellberg/hn"
)

func FetchLatest(count int) []*hn.Item {
	var items []*hn.Item

	hn := hn.DefaultClient

	ids, err := hn.TopStories()
	if err != nil {
		panic(err)
	}

	for _, id := range ids[:count] {
		item, err := hn.Item(id)
		if err != nil {
			panic(err)
		}

		items = append(items, item)
	}

	return items
}
