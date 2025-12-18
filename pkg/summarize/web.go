package summarize

import (
	"io"
	"log"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

func extractTextFromHTML(htmlContent string) string {
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		log.Printf("Error parsing HTML: %v", err)
		return ""
	}

	var text strings.Builder
	var extractText func(*html.Node)
	extractText = func(n *html.Node) {
		// Skip script, style, and other non-visible elements
		if n.Type == html.ElementNode {
			switch n.Data {
			case "script", "style", "noscript", "iframe", "object", "embed":
				return
			}
		}

		// Extract text from text nodes
		if n.Type == html.TextNode {
			content := strings.TrimSpace(n.Data)
			if content != "" {
				if text.Len() > 0 {
					text.WriteString(" ")
				}
				text.WriteString(content)
			}
		}

		// Recursively process child nodes
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extractText(c)
		}
	}

	extractText(doc)
	return strings.TrimSpace(text.String())
}

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
