package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/time/rate"

	"news-api/db"
	"news-api/handlers"
)

var RssSources = []string{
	"https://www.bleepingcomputer.com/feed/",
	"https://feeds.feedburner.com/TheHackersNews",
	"https://krebsonsecurity.com/feed/",
	"https://blogs.cisco.com/security/feed",
	"https://www.wired.com/feed/category/security/latest/rss",
	"https://www.securityweek.com/feed/",
}

// Create a rate limiter that allows 1 request per second with a burst size of 5.
var limiter = rate.NewLimiter(1, 5)

func main() {
	// Start the background caching job
	db.StartCachingJob(RssSources)

	// The main handler is now wrapped in our security middlewares.
	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("./test"))
	mux.Handle("/", fs)
	mux.HandleFunc("/news", handlers.GetNews)
	mux.HandleFunc("/ad", handlers.GetAd)
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Chain the middlewares. The request will flow from logging to security headers to the rate limiter.
	handler := loggingMiddleware(securityHeadersMiddleware(rateLimitMiddleware(mux)))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server starting on port " + port + "...")
	log.Fatal(http.ListenAndServe(":"+port, handler))
}

// Middleware for logging requests
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s %s", r.Method, r.RequestURI, r.RemoteAddr, time.Since(start))
	})
}

// Middleware to add security headers
func securityHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline';")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		next.ServeHTTP(w, r)
	})
}

// Middleware for rate limiting
func rateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
