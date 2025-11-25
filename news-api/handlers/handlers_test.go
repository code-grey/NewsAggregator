package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"news-api/db"
	"news-api/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTestDB initializes a clean in-memory database for testing.
func setupTestDB(t *testing.T) {
	if err := db.InitDB(":memory:"); err != nil {
		t.Fatalf("Failed to initialize test database: %v", err)
	}
}

// clearDB is a helper to clear articles between test runs.
func clearDB(t *testing.T) {
	err := db.ClearAllArticlesForTest()
	if err != nil {
		t.Fatalf("Failed to clear database: %v", err)
	}
}

// seedArticles inserts a set of predefined articles for testing filters.
func seedArticles(t *testing.T) {
	clearDB(t) // Ensure clean state before seeding
	now := time.Now()
	articles := []models.NewsArticle{
		{Title: "Cyber Article 1", URL: "u1", SourceURL: "src1", Category: "Cybersecurity", PublishedAt: now.Add(-1 * time.Hour), Rank: 10},
		{Title: "Tech Article 1", URL: "u2", SourceURL: "src2", Category: "Tech", PublishedAt: now.Add(-2 * time.Hour), Rank: 5},
		{Title: "Cyber Article 2 about ransomware", URL: "u3", SourceURL: "src1", Category: "Cybersecurity", PublishedAt: now.Add(-3 * time.Hour), Rank: 8},
		{Title: "Old Tech Article", URL: "u4", SourceURL: "src2", Category: "Tech", PublishedAt: now.Add(-48 * time.Hour), Rank: 2},
	}

	for _, article := range articles {
		err := db.InsertArticle(article)
		require.NoError(t, err)
	}
}

func TestGetNewsDefaultLimit(t *testing.T) {
	setupTestDB(t)
	clearDB(t)

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
		db.InsertArticle(article)
	}

	req, err := http.NewRequest("GET", "/news", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetNews)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "handler returned wrong status code")

	var responseArticles []models.NewsArticle
	err = json.NewDecoder(rr.Body).Decode(&responseArticles)
	require.NoError(t, err, "could not decode response")

	assert.Len(t, responseArticles, 20, "handler returned unexpected number of articles")
}

func TestGetNewsWithFilters(t *testing.T) {
	setupTestDB(t)
	seedArticles(t)

	testCases := []struct {
		name           string
		url            string
		expectedCount  int
		expectedTitles []string
	}{
		{
			name:          "Filter by category",
			url:           "/news?category=Cybersecurity",
			expectedCount: 2,
			expectedTitles: []string{"Cyber Article 1", "Cyber Article 2 about ransomware"},
		},
		{
			name:          "Filter by source",
			url:           "/news?source=src2",
			expectedCount: 2,
			expectedTitles: []string{"Tech Article 1", "Old Tech Article"},
		},
		{
			name:          "Filter by search term",
			url:           "/news?search=ransomware",
			expectedCount: 1,
			expectedTitles: []string{"Cyber Article 2 about ransomware"},
		},
		{
			name:          "Filter by date range",
			url:           "/news?start=" + time.Now().Add(-5*time.Hour).Format("2006-01-02") + "&end=" + time.Now().Format("2006-01-02"),
			expectedCount: 3, // Excludes the old article
		},
		{
			name:           "Sort by rank",
			url:            "/news?sortBy=rank",
			expectedCount:  4,
			expectedTitles: []string{"Cyber Article 1", "Cyber Article 2 about ransomware", "Tech Article 1", "Old Tech Article"},
		},
		{
			name:          "Limit results",
			url:           "/news?limit=1&sortBy=rank",
			expectedCount: 1,
			expectedTitles: []string{"Cyber Article 1"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", tc.url, nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(GetNews)
			handler.ServeHTTP(rr, req)

			assert.Equal(t, http.StatusOK, rr.Code)

			var articles []models.NewsArticle
			err = json.NewDecoder(rr.Body).Decode(&articles)
			require.NoError(t, err)

			assert.Len(t, articles, tc.expectedCount)

			if len(tc.expectedTitles) > 0 {
				var titles []string
				for _, a := range articles {
					titles = append(titles, a.Title)
				}
				assert.Equal(t, tc.expectedTitles, titles)
			}
		})
	}
}

func TestGetNewsInvalidDate(t *testing.T) {
	setupTestDB(t)

	req, err := http.NewRequest("GET", "/news?start=invalid-date", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetNews)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestGetTodayThreat(t *testing.T) {
	setupTestDB(t)
	seedArticles(t) // Seeds articles with various ranks and timestamps

	req, err := http.NewRequest("GET", "/today-threat", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetTodayThreat)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var threatScore db.ThreatScore
	err = json.NewDecoder(rr.Body).Decode(&threatScore)
	require.NoError(t, err)

	// Based on seedArticles (only recent articles count towards today's threat):
	// Recent: Cyber Article 1 (rank 10), Tech Article 1 (rank 5), Cyber Article 2 (rank 8)
	// High (rank >= 5): 3
	// Medium (2 <= rank < 5): 0
	// Low (rank < 2): 0
	// Total: 3
	assert.Equal(t, 3, threatScore.HighRankCount)
	assert.Equal(t, 0, threatScore.MediumRankCount)
	assert.Equal(t, 0, threatScore.LowRankCount)
	assert.Equal(t, 3, threatScore.TotalArticles)
	assert.Equal(t, "Code Red", threatScore.ThreatLevel)
}

func TestExportCSV(t *testing.T) {
	setupTestDB(t)
	seedArticles(t)

	req, err := http.NewRequest("GET", "/export/csv", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ExportCSV)
	handler.ServeHTTP(rr, req)

	// Check status and headers
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "text/csv", rr.Header().Get("Content-Type"))
	assert.Equal(t, `attachment; filename="articles.csv"`, rr.Header().Get("Content-Disposition"))

	// Check CSV content
	body := rr.Body.String()
	assert.Contains(t, body, "Title,Description,ImageURL,URL,SourceURL,PublishedAt,Rank,Category\n", "CSV header is missing or incorrect")
	assert.Contains(t, body, "Cyber Article 1,", "CSV should contain data from seeded articles")
	assert.Contains(t, body, "Tech Article 1,", "CSV should contain data from seeded articles")
}
