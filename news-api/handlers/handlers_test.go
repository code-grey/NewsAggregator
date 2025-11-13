package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"

	"news-api/db"
	"news-api/models"
)

func setupTestDB(t *testing.T) {
	// Use an in-memory SQLite database for testing
	if err := db.InitDB(); err != nil {
		t.Fatalf("Failed to initialize test database: %v", err)
	}

	// Add more than 20 articles for testing the default limit
	for i := 0; i < 25; i++ {
		article := models.NewsArticle{
			Title:       "Test Article " + strconv.Itoa(i),
			Description: "Test Description",
			URL:         "http://test.com/" + strconv.Itoa(i),
			SourceURL:   "http://testsource.com",
			PublishedAt: time.Now(),
			Category:    "Tech",
		}
		// Directly use the internal insertArticle function for testing
		// This is a simplified approach for the test setup.
		// In a real-world scenario, you might have a dedicated test helper for this.
		db.InsertArticle(article)
	}
}

func TestGetNewsDefaultLimit(t *testing.T) {
	setupTestDB(t)
	defer os.Remove("./news.db") // Clean up the database file

	req, err := http.NewRequest("GET", "/news", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetNews)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var articles []models.NewsArticle
	if err := json.NewDecoder(rr.Body).Decode(&articles); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}

	// This will fail before the fix, as it will return all 25 articles
	if len(articles) != 20 {
		t.Errorf("handler returned unexpected number of articles: got %v want %v",
			len(articles), 20)
	}
}
