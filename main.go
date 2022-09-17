package main

import (
	"fmt"
	"github.com/notarock/technews-bot/pkg/sources/hackernews"
)

func main() {
	i := hackernews.FetchLatest(10)
	fmt.Println(i)
}
