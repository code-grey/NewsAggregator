# Project Context: NewsAggregator Development Log

This document provides a detailed summary of the development process for the NewsAggregator project, covering all modifications, discussions, and decisions made.

## 1. Initial State & Problem Identification

The project started with a Go backend for news aggregation and a simple HTML frontend (`news-api/test/index.html`). The user reported two main issues:
1.  Problems with `@news-api/test/index.html` display.
2.  A new news source, "SecurityWeek," was not showing up in the dropdown menu.

## 2. Initial Fixes & RSS Feed Management

*   **SecurityWeek Integration:** Identified the RSS feed for SecurityWeek (`http://feeds.feedburner.com/Securityweek`) via a Google search.
*   **Frontend HTML Update (`news-api/test/index.html`):**
    *   Added the SecurityWeek RSS feed to the `rssSources` array in the JavaScript.
    *   Improved the display of articles by changing the `border-radius` from `9999px` (pill shape) to `4px` (card format).
    *   Enhanced ad display and error handling in the JavaScript `fetchAd` and `fetchNews` functions.
    *   Modified the domain extraction logic for RSS sources to handle `feedburner.com` URLs more gracefully.
*   **Backend Port Conflict:** The server failed to start due to port 8080 being in use.
    *   Initially attempted to change the port to 8081 in `news-api/main.go`.

## 3. Backend Restructuring (Go Modules)

A significant refactoring was undertaken to improve the Go backend's structure and maintainability, addressing build errors related to package imports.

*   **Problem:** `go build` failed with errors like "package news-api/handlers is not in std" and "replacement directory ./handlers does not exist." This indicated that the Go module system wasn't correctly recognizing local packages.
*   **Solution:**
    *   **Created Subdirectories:** Created `news-api/handlers`, `news-api/models`, and `news-api/db` directories.
    *   **Moved Files:**
        *   `handlers.go` moved to `news-api/handlers/handlers.go`.
        *   `models.go` moved to `news-api/models/models.go`.
        *   `db.go` moved to `news-api/db/db.go`.
    *   **Updated `go.mod`:** Added `replace` directives to `news-api/go.mod` to correctly map local package imports:
        ```go
        replace (
            "news-api/handlers" => "./handlers"
            "news-api/models" => "./models"
            "news-api/db" => "./db"
        )
        ```
    *   **Updated Package Declarations:** Changed `package main` to `package handlers`, `package models`, and `package db` in their respective files.
    *   **Updated Import Paths:** Modified import statements in `main.go`, `handlers/handlers.go`, and `db/db.go` to use the new module-relative paths (e.g., `"news-api/handlers"`, `"news-api/models"`, `"news-api/db"`).
    *   **Public Functions:** Capitalized the first letter of `InitDB`, `GetArticlesFromDB`, and `StartCachingJob` in `db/db.go` to make them public and accessible from other packages.
    *   **`rssSources` Visibility:** Made `rssSources` in `main.go` public (`RssSources`) and passed it as a parameter to `db.StartCachingJob`.

## 4. Addressing RSS Feed Issues (Dark Reading & Threatpost)

*   **Dark Reading:** The `https://www.darkreading.com/rss/all.xml` feed was identified as problematic. After further investigation, it was decided to remove it entirely due to persistent issues.
*   **Threatpost:** The `threatpost.com/feed/` URL was found to be a dead link. It was replaced with `https://www.wired.com/feed/category/security/latest/rss` (Wired Security feed).
*   **Updates:** Removed/updated these feeds in both `news-api/main.go` and `news-api/test/index.html`.

## 5. Implementing Keyword-Based Ranking System

A new feature was added to rank news articles based on keyword presence.

*   **`models/models.go`:** Added a `Rank` field (`int`) to the `NewsArticle` struct.
*   **`db/db.go`:**
    *   Modified the `articles` table creation SQL to include a `rank INTEGER DEFAULT 0` column.
    *   Implemented a `calculateRank` function that assigns scores based on high, medium, and low-impact cybersecurity keywords found in the article's title and description.
    *   Updated `insertArticle` to store the calculated rank.
    *   Modified `GetArticlesFromDB` to accept a `sortBy` string parameter. If `sortBy` is "rank", it orders results by `rank DESC`; otherwise, it defaults to `publishedAt DESC`.
*   **`handlers/handlers.go`:** Modified the `GetNews` handler to retrieve the `sortBy` query parameter from the request and pass it to `db.GetArticlesFromDB`.
*   **`news-api/test/index.html`:**
    *   Added a "Sort by" dropdown (`<select id="sort-by">`) with options for "Date" (default) and "Rank".
    *   Updated the `fetchNews` JavaScript function to include the `sortBy` parameter in the API request.
    *   Modified the article display to show the calculated rank.
*   **Database Schema Migration:** To apply the new `rank` column to the existing database, the old `news.db` file was deleted. This forces the application to create a new database with the updated schema on startup.

## 6. Deployment Preparation for Render

Key adjustments were made to ensure smooth deployment on Render's free tier.

*   **Dynamic Port Binding:** Modified `news-api/main.go` to read the listening port from the `PORT` environment variable (provided by Render). It falls back to `8080` if the variable is not set (for local development).
*   **Health Check Endpoint:** Added a `/healthz` endpoint to `main.go` that returns a `200 OK` status. This is crucial for external cron jobs to keep the Render service alive and prevent it from sleeping due to inactivity.
*   **Cross-Platform Binaries:** Compiled the Go backend for various platforms:
    *   Windows (AMD64): `news-api-windows-amd64.exe`
    *   Linux (AMD64): `news-api-linux-amd64`
    *   macOS (AMD64): `news-api-darwin-amd64`
    *   macOS (ARM64): `news-api-darwin-arm64`
    *   Note: Required using `set GOOS=... && set GOARCH=... && go build` syntax on Windows.

## 7. Git Management

*   **`.gitignore`:** Added `news.db` to `.gitignore` to prevent it from being tracked by Git.
*   **Commits:** All changes were staged and committed with descriptive messages, using a temporary file for the commit message to avoid shell quoting issues.

## 8. Documentation

*   **`architecture.md`:** Created a comprehensive document detailing the project's overall architecture, backend components, features, data flow, and deployment considerations.
*   **`optimize.md`:** Created a document outlining various Android app optimization techniques, focusing on reducing app size (App Bundles, R8, dynamic features, vector drawables, image compression, lightweight libraries, lazy loading).

## 9. Discussions & Future Considerations

Throughout the development, several important architectural and feature discussions took place:

*   **Article Summarization:**
    *   **Challenge:** RSS feeds often provide only snippets, not full articles.
    *   **Solution:** Scraping the full article content from source URLs using `goquery` (for extraction) and `bluemonday` (for sanitization).
    *   **Website Scraping Policies:** Investigated `robots.txt` files for all RSS sources to determine scraping allowances (e.g., Krebs on Security disallows).
*   **Semantic Similarity Search:**
    *   **Feasibility:** Confirmed it's possible with small language models (<1GB).
    *   **Implementation:** Generate sentence embeddings for articles, store them, and use cosine similarity for comparison.
    *   **Model Hosting:** Recommended **client-side (on-device)** hosting for the embedding model (e.g., Gemini Nano) due to Render's free-tier resource limitations and cost implications. This offloads processing from the server.
*   **Client-Side SQLite for Archive:**
    *   **Proposal:** Store a persistent archive of news articles directly on the user's Android device using SQLite.
    *   **Benefits:** Offline access, user-specific history, improved performance, and keeps the backend stateless and free-tier compatible.
    *   **Cleanup:** Implement a weekly cleanup mechanism for older articles in the client-side DB.
*   **Cold Storage Archive:**
    *   **Discussion:** Explored options for a permanent archive of article URLs.
    *   **Conclusion:** Client-side SQLite is the most practical for an MVP. Server-side CSV/Google Docs were deemed impractical for the free tier due to ephemeral storage or excessive complexity.
*   **Fetching Frequency:** The server fetches news every 15 minutes. The client fetches news from the server on user demand (app open, refresh).
*   **Render Free Tier Limitations:** Acknowledged and planned around ephemeral storage and cold starts. The health check endpoint and external cron job are solutions for the latter.

This comprehensive log should provide all the necessary context to resume development at any point.
