package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"news-api/db"
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