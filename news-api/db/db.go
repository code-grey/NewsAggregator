package db

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/microcosm-cc/bluemonday"
	"github.com/mmcdole/gofeed"
	"news-api/models"
	"github.com/pemistahl/lingua-go"
)

var db *sql.DB
var detector lingua.LanguageDetector

func InitDB() error {
	var err error
	db, err = sql.Open("sqlite3", "./news.db")
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}

	createTableSQL := `
	CREATE TABLE IF NOT EXISTS articles (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		description TEXT,
		imageUrl TEXT,
		url TEXT NOT NULL UNIQUE,
		sourceUrl TEXT NOT NULL,
		publishedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
		rank INTEGER DEFAULT 0
	);
	`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("failed to create articles table: %v", err)
	}

	// Create indexes for faster queries
	createIndexesSQL := `
	CREATE INDEX IF NOT EXISTS idx_sourceUrl ON articles (sourceUrl);
	CREATE INDEX IF NOT EXISTS idx_publishedAt ON articles (publishedAt);
	`
	_, err = db.Exec(createIndexesSQL)
	if err != nil {
		return fmt.Errorf("failed to create indexes: %v", err)
	}

	// Optimize language detector to only load models for relevant languages
	detector = lingua.NewLanguageDetectorBuilder().
		FromLanguages(lingua.English, lingua.German, lingua.French, lingua.Spanish, lingua.Russian, lingua.Chinese).
		WithPreloadedLanguageModels().
		Build()

	log.Println("Database initialized successfully.")
	return nil
}

func calculateRank(article models.NewsArticle) int {
	rank := 0
	// Refined keywords focusing on active threats and their severity
	keywords := map[string]int{
		// High Impact (Score 5): Direct, immediate threats
		"zero-day": 5, "exploit in the wild": 5, "active attack": 5, "critical vulnerability": 5, "alert": 5, "warning": 5, "patch now": 5, "ransomware attack": 5, "breach confirmed": 5,
		// Medium Impact (Score 3): Significant threats, but perhaps not immediate action required
		"vulnerability": 3, "exploit": 3, "breach": 3, "attack": 3, "malware": 3, "ransomware": 3, "phishing": 3, "threat": 3, "advisory": 3,
		// Low Impact (Score 1): General cybersecurity news, informative
		"security": 1, "cybersecurity": 1, "data": 1, "privacy": 1, "risk": 1, "compliance": 1, "encryption": 1, "patch": 1,
	}

	content := strings.ToLower(article.Title + " " + article.Description)

	for keyword, score := range keywords {
		if strings.Contains(content, keyword) {
			rank += score
		}
	}

	return rank
}

func insertArticle(article models.NewsArticle) error {
	stmt, err := db.Prepare("INSERT OR IGNORE INTO articles(title, description, imageUrl, url, sourceUrl, publishedAt, rank) VALUES(?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Printf("Error preparing insert statement for article %s: %v", article.Title, err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(article.Title, article.Description, article.ImageURL, article.URL, article.SourceURL, article.PublishedAt, article.Rank)
	if err != nil {
		log.Printf("Error inserting article %s: %v", article.Title, err)
	}
	return err
}

// ThreatScore represents the calculated threat score and its corresponding phrase.
type ThreatScore struct {
	Score  float64 `json:"score"`
	Phrase string  `json:"phrase"`
}

const MIN_ARTICLES_FOR_SCORE = 5 // Minimum articles required for a reliable score

// GetTodayThreatScore calculates the average rank of articles published in the last 24 hours.
func GetTodayThreatScore() (ThreatScore, error) {
	var totalRank int
	var articleCount int

	// Define the time window (last 24 hours)
	twentyFourHoursAgo := time.Now().Add(-24 * time.Hour)

	rows, err := db.Query("SELECT rank FROM articles WHERE publishedAt >= ?", twentyFourHoursAgo.Format("2006-01-02 15:04:05"))
	if err != nil {
		return ThreatScore{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var rank int
		if err := rows.Scan(&rank); err != nil {
			log.Printf("Error scanning rank for threat score: %v", err)
			continue
		}
		totalRank += rank
		articleCount++
	}

	if articleCount < MIN_ARTICLES_FOR_SCORE {
		return ThreatScore{Score: 0, Phrase: "No Worries (Insufficient Data)"}, nil
	}

	averageRank := float64(totalRank) / float64(articleCount)

	phrase := "No Worries" // Default
	if averageRank >= 1.6 && averageRank <= 3.5 {
		phrase = "Attention!"
	} else if averageRank > 3.5 {
		phrase = "Code Red"
	}

	return ThreatScore{Score: averageRank, Phrase: phrase}, nil
}

func GetArticlesFromDB(sourceFilter string, limit int, startDate, endDate time.Time, sortBy string) ([]models.NewsArticle, error) {
	if db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}
	var articles []models.NewsArticle
	query := "SELECT title, description, imageUrl, url, sourceUrl, publishedAt, rank FROM articles"
	args := []interface{}{}

	whereClauses := []string{}

	if sourceFilter != "" && sourceFilter != "all" {
		whereClauses = append(whereClauses, "sourceUrl = ?")
		args = append(args, sourceFilter)
	}

	if !startDate.IsZero() {
		whereClauses = append(whereClauses, "publishedAt >= ?")
		args = append(args, startDate.Format("2006-01-02 15:04:05"))
	}
	if !endDate.IsZero() {
		whereClauses = append(whereClauses, "publishedAt <= ?")
		args = append(args, endDate.Format("2006-01-02 15:04:05"))
	}

	if len(whereClauses) > 0 {
		query += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	if sortBy == "rank" {
		query += " ORDER BY rank DESC"
	} else {
		query += " ORDER BY publishedAt DESC"
	}

	if limit > 0 {
		query += " LIMIT ?"
		args = append(args, limit)
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		log.Printf("Error executing query in GetArticlesFromDB: %v", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var article models.NewsArticle
		if err := rows.Scan(&article.Title, &article.Description, &article.ImageURL, &article.URL, &article.SourceURL, &article.PublishedAt, &article.Rank); err != nil {
			log.Printf("Error scanning article: %v", err)
			continue
		}
		articles = append(articles, article)
	}

	return articles, nil
}

func StartCachingJob(rssSources []string) {
	if err := InitDB(); err != nil {
		log.Printf("Failed to initialize database for caching job: %v", err)
		return
	}
	fetchAndCacheNews(rssSources)

	ticker := time.NewTicker(15 * time.Minute)
	go func() {
		for range ticker.C {
			log.Println("Running scheduled news caching job...")
			fetchAndCacheNews(rssSources)
		}
	}()
}

func fetchAndCacheNews(rssSources []string) {
	client := &http.Client{Timeout: 10 * time.Second}
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout: 10 * time.Second,
	}
	client.Transport = &userAgentTransport{RoundTripper: transport}

	fp := gofeed.NewParser()
	fp.Client = client

	var wg sync.WaitGroup
	p := bluemonday.StripTagsPolicy()

	for _, source := range rssSources {
		wg.Add(1)
		go func(source string) {
			defer wg.Done()
			feed, err := fp.ParseURL(source)
			if err != nil {
				log.Printf("Error parsing feed from %s for caching: %v", source, err)
				return
			}

			for _, item := range feed.Items {
				// Language detection
				textToDetect := item.Title + " " + item.Description
				lang, _ := detector.DetectLanguageOf(textToDetect)
				if lang != lingua.English {
					log.Printf("Skipping non-English article: %s (Source: %s)", item.Title, source)
					continue
				}

				article := models.NewsArticle{
					Title:       item.Title,
					Description: p.Sanitize(item.Description),
					URL:         item.Link,
					SourceURL:   source,
				}
				article.Rank = calculateRank(article)

				if item.Image != nil {
					article.ImageURL = item.Image.URL
				}
				if item.PublishedParsed != nil {
					article.PublishedAt = *item.PublishedParsed
				} else if feed.PublishedParsed != nil {
					article.PublishedAt = *feed.PublishedParsed
				} else {
					article.PublishedAt = time.Now()
				}

				if err := insertArticle(article); err != nil {
					// log.Printf("Error inserting article %s: %v", article.Title, err) // Log only if not a unique constraint violation
				}
			}
		}(source)
	}

	wg.Wait()
	log.Println("News caching job completed.")
}

type userAgentTransport struct {
	http.RoundTripper
}

func (t *userAgentTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36")
	return t.RoundTripper.RoundTrip(req)
}