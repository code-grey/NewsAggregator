package db

import (
	"database/sql"
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
)

var db *sql.DB

func InitDB() {
	var err error
	db, err = sql.Open("sqlite3", "./news.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
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
		log.Fatalf("Failed to create articles table: %v", err)
	}
	log.Println("Database initialized successfully.")
}

func calculateRank(article models.NewsArticle) int {
	rank := 0
	keywords := map[string]int{
		"vulnerability": 3, "exploit": 3, "breach": 3, "attack": 3, "malware": 3, "ransomware": 3,
		"security": 2, "cybersecurity": 2, "threat": 2, "phishing": 2, "patch": 2, "zero-day": 2,
		"data": 1, "privacy": 1, "risk": 1, "compliance": 1, "encryption": 1,
	}

	for keyword, score := range keywords {
		if strings.Contains(strings.ToLower(article.Title), keyword) || strings.Contains(strings.ToLower(article.Description), keyword) {
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

func GetArticlesFromDB(sourceFilter string, limit int, startDate, endDate time.Time, sortBy string) ([]models.NewsArticle, error) {
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
	InitDB()
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

