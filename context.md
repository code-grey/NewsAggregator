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
    *   Enhanced error handling in the JavaScript `fetchNews` function.
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

*   **`models.go`:** Added a `Rank` field (`int`) to the `NewsArticle` struct.
*   **`db.go`:**
    *   Modified the `articles` table creation SQL to include a `rank INTEGER DEFAULT 0` column.
    *   Implemented a `calculateRank` function that assigns scores based on high, medium, and low-impact cybersecurity keywords found in the article's title and description.
    *   Updated `insertArticle` to store the calculated rank.
    *   Modified `GetArticlesFromDB` to accept a `sortBy` string parameter. If `sortBy` is "rank", it orders results by `rank DESC`; otherwise, it defaults to `publishedAt DESC`.
*   **`handlers.go`:** Modified the `GetNews` handler to retrieve the `sortBy` query parameter from the request and pass it to `db.GetArticlesFromDB`.
*   **`index.html`:**
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

## 9. Recent Developments & Feature Implementations

This section details the work done since the last major update, focusing on stabilizing the backend, enhancing the user experience, and introducing new features.

### 9.1 Backend Stability & Debugging

*   **`git reset --hard` Mishap & Re-implementation:** An accidental `git reset --hard` command led to the loss of uncommitted changes, necessitating the re-implementation of several features and fixes. This included re-applying changes to `main.go`, `db.go`, `handlers.go`, and `index.html` to restore the intended functionality.

*   **Persistent `500 Internal Server Error`:** Encountered persistent `500 Internal Server Error` when fetching articles, despite successful database initialization and caching.
    *   **Diagnosis:** Initial attempts to debug by adding detailed error logging to `handlers.go` and `db.go` were unsuccessful in revealing the root cause directly in the browser.
    *   **Root Cause Identification:** Suspected issues with the `db` connection being `nil` or invalid, or `log.Fatalf` calls terminating the server prematurely.
    *   **Solution:**
        *   Modified `db.InitDB()` to return an error instead of calling `log.Fatalf`.
        *   Modified `main.go` to handle the error returned by `db.InitDB()` using `log.Fatalf`, ensuring that database initialization failures are explicitly logged and stop the server.
        *   Added a `nil` check for the `db` connection in `db.GetArticlesFromDB` to return a more explicit error if the connection is not established.
        *   Removed the temporary `check_db.go` file to prevent build conflicts.
*   **Server Accessibility:** Ensured the server is running in the background using `start /B` for continuous operation during development.

### 9.2 Frontend Enhancements & Bug Fixes

*   **Article Display Issues:**
    *   **Problem:** Articles were not showing up on the frontend despite successful backend caching.
    *   **Root Cause:** The `http.FileServer` in `main.go` was incorrectly capturing all requests due to `mux.Handle("/", fs)`, preventing API endpoints from being reached.
    *   **Solution:** Changed `mux.Handle("/", fs)` to `mux.Handle("/static/", http.StripPrefix("/static/", fs))` in `main.go` to serve static files only from the `/static/` path.
*   **Image Removal:** Removed the display of images from articles in `news-api/test/index.html` for a cleaner, text-focused view.
*   **Instant Filter Application:** Implemented `change` event listeners on filter `select` and `input[type="date"]` elements in `news-api/test/index.html` to trigger `fetchAll()` instantly, removing the need for an "Apply" button.
*   **Duplicate Rank Display:** Fixed a bug in `news-api/test/index.html` that caused the article rank to be displayed twice.
*   **CSO Online Visibility:** Re-added `https://www.csoonline.com/feed/` to both `news-api/main.go` and `news-api/test/index.html` to make it visible in the source dropdown, acknowledging its mixed-language content for now.

### 9.3 "Today in ThreatFeed" Feature Implementation (Deprecated)

*   **Concept:** Introduced a new feature to provide an at-a-glance summary of the current cybersecurity threat level.
*   **Scoring Mechanism:**
    *   **Initial Idea:** Simple sum of article ranks.
    *   **Refinement:** Adopted an **average rank** of articles published in the last 24 hours to provide a more meaningful "per-article" severity.
    *   **Thresholds:** Defined numerical thresholds for the average rank to map to qualitative threat levels:
        *   `0.0` to `1.5`: "No Worries" (Green)
        *   `1.6` to `3.5`: "Attention!" (Yellow/Orange)
        *   `3.6` and above: "Code Red" (Red)
    *   **Insufficient Data Handling:** Implemented a `MIN_ARTICLES_FOR_SCORE` constant (e.g., 5 articles) to return "No Worries (Insufficient Data)" if not enough articles are available for a reliable score.
*   **Backend Implementation:**
    *   Added `ThreatScore` struct and `GetTodayThreatScore()` function in `news-api/db/db.go` to calculate the average rank and phrase.
    *   Created `/today-threat` API endpoint in `news-api/handlers/handlers.go` to expose the threat score.
*   **Frontend Implementation:**
    *   Added a new HTML section (`#threat-score-section`) in `news-api/test/index.html` to display the score and a meter bar.
    *   Added CSS for `threat-meter-container` and `threat-meter-fill` with color classes (`threat-low`, `threat-medium`, `threat-high`) for visual representation.
    *   Implemented JavaScript to fetch the score from `/today-threat`, update the phrase, and dynamically control the meter's width and color.

### 9.4 Future Considerations & Discussions

*   **"What This Means For You" Section:** Discussed the potential for a highly valuable feature that translates aggregated threat data into personalized, actionable implications for the user.
    *   **Pros:** High value, personalization, proactive security, strong differentiation.
    *   **Cons:** Significant technical complexity (especially for personalized advice requiring NLP and knowledge bases), risk of misinformation, high maintenance.
    *   **Confidence Score:** 7/10 for a basic version (general advice), 3/10 for a highly personalized version without additional resources.
    *   **Proposed Phased Approach:** Start with general advice (V1), then explore curated explanations (V2), and finally advanced PoC linking (V3) with careful risk management.
*   **Language Detection (Next Step):** Identified the need to implement language detection in the Go backend to filter out non-English articles from mixed-language feeds, ensuring data quality. Lingua-Go was identified as a suitable library.
*   **Vulnerability to PoC/Explanation Linkage:** Explored the idea of linking vulnerabilities to PoC exploits, MITRE CVEs, and technical explanations. This was deemed a highly valuable but complex feature with significant safety and maintenance considerations, recommending a phased implementation.

### 9.5 ThreatFeed Score Redesign

*   **Concept:** Replaced the single average rank and phrase with a visual representation of the percentage breakdown of high, medium, and low-ranked articles for the day.
*   **Backend Implementation:**
    *   Modified `ThreatScore` struct to return `LowRankCount`, `MediumRankCount`, `HighRankCount`, and `TotalArticles`.
    *   Updated `GetTodayThreatScore()` to calculate and return these counts.
*   **Frontend Implementation:**
    *   Modified HTML structure of `#threat-score-section` to accommodate multiple colored bars.
    *   Updated JavaScript in `fetchThreatScore()` to parse the new API response and dynamically create/style percentage bars.
    *   Added `threat-bar` CSS class for styling.

### 9.6 New Content Categories (Tech & Defense)

*   **Concept:** Expanded content categories beyond Cybersecurity to include "Tech" and "Defense" news.
*   **Backend Implementation:**
    *   Updated `NewsArticle` model and database schema with a `category` field.
    *   Integrated new RSS feeds for "Tech" and "Defense" news.
    *   Modified `calculateRank` to use category-specific keywords for ranking.
    *   Updated `getCategoryForSource` to map new feeds to their respective categories.
    *   Enabled category filtering in the backend API.
*   **Frontend Implementation:**
    *   Added "Tech" and "Defense" options to the category dropdown.
    *   Updated `rssSources` array to include new feeds.
    *   Implemented a loading indicator for article fetching.

This comprehensive log should provide all the necessary context to resume development at any point.
