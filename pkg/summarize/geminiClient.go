package summarize

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/genai"
)

type GeminiClient struct {
	client *genai.Client
}

func InitGeminiClient() (*GeminiClient, error) {
	ctx := context.Background()
	// The client gets the API key from the environment variable `GEMINI_API_KEY`.
	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	return &GeminiClient{client: client}, nil
}

func (gc *GeminiClient) SummarizeWebpage(URL string) (string, error) {
	// fetch the html page
	html, err := fetchWebpage(URL)

	if err != nil {
		return "", fmt.Errorf("failed to fetch webpage: %w", err)
	}

	summary, err := gc.summarizeFromHtml(html)
	if err != nil {
		return "", fmt.Errorf("failed to summarize HTML: %w", err)
	}

	return summary, nil
}

func (gc *GeminiClient) summarizeFromHtml(htmlContent string) (string, error) {
	ctx := context.Background()

	prompt := fmt.Sprintf("Summarize the following content in no longer than 8 sentences:\n\n%s", htmlContent)

	result, err := gc.client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash",
		genai.Text(prompt),
		nil,
	)

	if err != nil {
		return "", err
	}

	return result.Text(), nil
}
