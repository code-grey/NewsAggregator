package handlers

import (
	"encoding/csv"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"news-api/db"
	"news-api/models"
)

func GetNews(w http.ResponseWriter, r *http.Request) {
	// Get query parameters
	sourceFilter := r.URL.Query().Get("source")
	categoryFilter := r.URL.Query().Get("category") // New parameter
	searchFilter := r.URL.Query().Get("search")
	limitStr := r.URL.Query().Get("limit")
	limit, _ := strconv.Atoi(limitStr)
	if limit == 0 {
		limit = 20 // Default limit
	}
	startDateStr := r.URL.Query().Get("start")
	endDateStr := r.URL.Query().Get("end")
	sortBy := r.URL.Query().Get("sortBy")

	var startDate, endDate time.Time
	var err error

	if startDateStr != "" {
		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			http.Error(w, "Invalid start date format", http.StatusBadRequest)
			return
		}
	}

	if endDateStr != "" {
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			http.Error(w, "Invalid end date format", http.StatusBadRequest)
			return
		}
		// Add 23 hours, 59 minutes, and 59 seconds to the end date to include the entire day.
		endDate = endDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
	}

	articles, err := db.GetArticlesFromDB(sourceFilter, categoryFilter, searchFilter, limit, startDate, endDate, sortBy) // Pass categoryFilter
	if err != nil {
		log.Printf("Error fetching articles from DB: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(articles)
}


func GetTodayThreat(w http.ResponseWriter, r *http.Request) {
	threatScore, err := db.GetTodayThreatScore()
	if err != nil {
		log.Printf("Error getting today's threat score: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(threatScore)
}

func ExportCSV(w http.ResponseWriter, r *http.Request) {
	// Set headers to prompt for file download.
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", `attachment; filename="articles.csv"`)

	rows, err := db.GetAllArticlesStream()
	if err != nil {
		log.Printf("Error getting articles stream from DB: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	csvWriter := csv.NewWriter(w)
	defer csvWriter.Flush()

	// Write CSV header
	headers := []string{"Title", "Description", "ImageURL", "URL", "SourceURL", "PublishedAt", "Rank", "Category"}
	if err := csvWriter.Write(headers); err != nil {
		log.Printf("Error writing CSV header: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Write rows
	for rows.Next() {
		var article models.NewsArticle
		if err := rows.Scan(&article.Title, &article.Description, &article.ImageURL, &article.URL, &article.SourceURL, &article.PublishedAt, &article.Rank, &article.Category); err != nil {
			log.Printf("Error scanning article row for CSV export: %v", err)
			continue // Skip bad rows
		}

		record := []string{
			article.Title,
			article.Description,
			article.ImageURL,
			article.URL,
			article.SourceURL,
			article.PublishedAt.Format(time.RFC3339), // Use a standard format
			strconv.Itoa(article.Rank),
			article.Category,
		}
		if err := csvWriter.Write(record); err != nil {
			log.Printf("Error writing CSV record: %v", err)
			// The connection might be broken, so we can't send another HTTP error.
			// We just log and stop.
			return
		}
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating article rows for CSV export: %v", err)
	}
}
