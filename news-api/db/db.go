package db

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"news-api/gemini"
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
		summary TEXT,
		imageUrl TEXT,
		url TEXT NOT NULL UNIQUE,
		sourceUrl TEXT NOT NULL,
		publishedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
		rank INTEGER DEFAULT 0,
		category TEXT DEFAULT ''
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
	content := strings.ToLower(article.Title + " " + article.Description)

	var keywords map[string]int

	switch article.Category {
	case "Cybersecurity":
		keywords = map[string]int{
			// High Impact (Score 5): Direct, immediate threats
			"zero-day": 5, "exploit in the wild": 5, "active attack": 5, "critical vulnerability": 5, "alert": 5, "warning": 5, "patch now": 5, "ransomware attack": 5, "breach confirmed": 5,
			// Medium Impact (Score 3): Significant threats, but perhaps not immediate action required
			"vulnerability": 3, "exploit": 3, "breach": 3, "attack": 3, "malware": 3, "ransomware": 3, "phishing": 3, "threat": 3, "advisory": 3,
			// Low Impact (Score 1): General cybersecurity news, informative
			"security": 1, "cybersecurity": 1, "data": 1, "privacy": 1, "risk": 1, "compliance": 1, "encryption": 1, "patch": 1,
		}
	case "Tech":
		keywords = map[string]int{
			// High Impact (Score 5): Major announcements, breakthroughs, critical issues
			"ai": 5, "artificial intelligence": 5, "quantum computing": 5, "breakthrough": 5, "major update": 5, "new chip": 5, "innovation": 5, "future of tech": 5,
			// Medium Impact (Score 3): Significant developments, new products, industry trends
			"startup": 3, "funding": 3, "acquisition": 3, "cloud": 3, "5g": 3, "machine learning": 3, "data science": 3, "web3": 3, "metaverse": 3, "robotics": 3,
			// Low Impact (Score 1): General tech news, reviews, minor updates
			"review": 1, "gadget": 1, "app": 1, "software": 1, "hardware": 1, "update": 1, "guide": 1, "tips": 1,
		}
	default: // General or unknown category
		keywords = map[string]int{
			"news": 1, "update": 1, "report": 1,
		}
	}

	for keyword, score := range keywords {
		if strings.Contains(content, keyword) {
			rank += score
		}
	}

	return rank
}

func insertArticle(article models.NewsArticle) error {
	stmt, err := db.Prepare("INSERT OR IGNORE INTO articles(title, description, summary, imageUrl, url, sourceUrl, publishedAt, rank, category) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Printf("Error preparing insert statement for article %s: %v", article.Title, err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(article.Title, article.Description, article.Summary, article.ImageURL, article.URL, article.SourceURL, article.PublishedAt, article.Rank, article.Category)
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

func GetArticlesFromDB(sourceFilter string, categoryFilter string, limit int, startDate, endDate time.Time, sortBy string) ([]models.NewsArticle, error) {
	if db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}
	var articles []models.NewsArticle
	query := "SELECT title, description, summary, imageUrl, url, sourceUrl, publishedAt, rank, category FROM articles"
	args := []interface{}{}

	whereClauses := []string{}

	if sourceFilter != "" && sourceFilter != "all" {
		whereClauses = append(whereClauses, "sourceUrl = ?")
		args = append(args, sourceFilter)
	}

	if categoryFilter != "" && categoryFilter != "all" {
		whereClauses = append(whereClauses, "category = ?")
		args = append(args, categoryFilter)
	}

	if !startDate.IsZero() {
		whereClauses = append(whereClauses, "publishedAt >= ?")
		args = append(args, startDate.Format("2006-01-02"))
	}
	if !endDate.IsZero() {
		whereClauses = append(whereClauses, "publishedAt <= ?")
		args = append(args, endDate.Format("2006-01-02"))
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
		if err := rows.Scan(&article.Title, &article.Description, &article.Summary, &article.ImageURL, &article.URL, &article.SourceURL, &article.PublishedAt, &article.Rank, &article.Category); err != nil {
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

				category := getCategoryForSource(source)

				sanitizedDescription := p.Sanitize(item.Description)
				var summary string
				apiKey := gemini.GetAPIKey()
				if apiKey != "" {
					var err error
					summary, err = gemini.SummarizeArticle(apiKey, sanitizedDescription)
					if err != nil {
						log.Printf("Error summarizing article with Gemini: %v. Falling back to truncation.", err)
						if len(sanitizedDescription) > 150 {
							summary = sanitizedDescription[:150] + "..."
						} else {
							summary = sanitizedDescription
						}
					}
				} else {
					if len(sanitizedDescription) > 150 {
						summary = sanitizedDescription[:150] + "..."
					} else {
						summary = sanitizedDescription
					}
				}

				article := models.NewsArticle{
					Title:       item.Title,
					Description: sanitizedDescription,
					Summary:     summary,
					URL:         item.Link,
					SourceURL:   source,
					Category:    category,
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

func getCategoryForSource(sourceURL string) string {
	// Define your source-to-category mapping here
	cybersecuritySources := []string{
		"https://www.bleepingcomputer.com/feed/",
		"https://feeds.feedburner.com/TheHackersNews",
		"https://blogs.cisco.com/security/feed",
		"https://www.wired.com/feed/category/security/latest/rss",
		"https://www.securityweek.com/feed/",
		"https://news.sophos.com/en-us/feed/",
		"https://www.csoonline.com/feed/",
	}

	techSources := []string{
		"https://www.theverge.com/rss/index.xml",
		"https://techcrunch.com/feed/",
		"https://arstechnica.com/feed/",
		"http://www.engadget.com/rss-full.xml",
		"http://www.fastcodesign.com/rss.xml",
		"http://www.forbes.com/entrepreneurs/index.xml",
		"https://blog.pragmaticengineer.com/rss/",
		"https://browser.engineering/rss.xml",
		"https://githubengineering.com/atom.xml",
		"https://joshwcomeau.com/rss.xml",
		"https://jvns.ca/atom.xml",
		"https://overreacted.io/rss.xml",
		"https://signal.org/blog/rss.xml",
		"https://slack.engineering/feed",
		"https://shopifyengineering.myshopify.com/blogs/engineering.atom",
		"https://stripe.com/blog/feed.rss",
		"https://www.uber.com/blog/engineering/rss/",
	}

	for _, s := range cybersecuritySources {
		if s == sourceURL {
			return "Cybersecurity"
		}
	}

	for _, s := range techSources {
		if s == sourceURL {
			return "Tech"
		}
	}

	return "General" // Default category if no match
}