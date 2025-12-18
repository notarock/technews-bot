package summarize

import (
	"io"
	"log"
	"net/http"
)

func fetchWebpage(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	// Check if the response status code indicates success.
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("unexpected status code: %d", resp.StatusCode)
	}

	// Read the entire response body into a byte slice.
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Convert the byte slice to a string and returns it.
	return string(bodyBytes), nil
}
