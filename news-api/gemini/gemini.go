package gemini

import (
	"context"
	"fmt"
	"os"

	"google.golang.org/genai"
)

// GetAPIKey retrieves the Gemini API key from the environment.
func GetAPIKey() string {
	return os.Getenv("GEMINI_API_KEY")
}

// SummarizeArticle uses the Gemini API to summarize the provided article content.
func SummarizeArticle(apiKey, articleContent string) (string, error) {
	if apiKey == "" {
		return "", fmt.Errorf("Gemini API key not provided")
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, genai.WithAPIKey(apiKey))
	if err != nil {
		return "", fmt.Errorf("failed to create genai client: %w", err)
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-1.5-flash")
	prompt := fmt.Sprintf("Summarize the following news article in 1-2 sentences:\n\n%s", articleContent)

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))

	if err != nil {
		return "", fmt.Errorf("failed to generate content: %w", err)
	}

	if len(resp.Candidates) > 0 {
		for _, part := range resp.Candidates[0].Content.Parts {
			if txt, ok := part.(genai.Text); ok {
				return string(txt), nil
			}
		}
	}

	return "", fmt.Errorf("no text content found in response")
}