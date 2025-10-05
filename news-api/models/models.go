package models

import "time"

// NewsArticle defines the structure for a news article.
type NewsArticle struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Summary     string `json:"summary"`
	ImageURL    string `json:"imageUrl"`
	URL         string `json:"url"`
	SourceURL   string `json:"sourceUrl"`
	PublishedAt time.Time `json:"publishedAt"`
	Rank        int    `json:"rank"`
	Category    string `json:"category"`
}

// Ad defines the structure for an ad.
type Ad struct {
	ImageURL  string `json:"imageUrl"`
	TargetURL string `json:"targetUrl"`
}