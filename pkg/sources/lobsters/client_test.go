package lobsters_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/notarock/technews-bot/pkg/sources/lobsters"
)

func TestFetchLatest(t *testing.T) {
	articles := lobsters.FetchLatest()

	assert.Equal(t, len(articles), 25)

	// Articles results are time dependant.
	// Lets just check we got something.
	assert.NotEmpty(t, articles[0].Title)
	assert.NotEmpty(t, articles[0].Link)
	assert.NotEmpty(t, articles[0].Tags)
}
