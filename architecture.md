# Project Architecture: NewsAggregator

This document outlines the architecture of the NewsAggregator project, a system designed to fetch, aggregate, and display cybersecurity news from various RSS sources.

## 1. Overall Architecture

The project follows a client-server architecture, with a Go-based backend responsible for news aggregation and API serving, and a frontend (currently a simple web interface, with an Android application planned) for user interaction.

```
+-----------------+       +-----------------+       +-----------------+
|   RSS Sources   | <---> |     Backend     | <---> |     Frontend    |
| (Bleeping Comp, |       |    (Go Lang)    |       | (Web / Android) |
|   The Hacker    |       |                 |       |                 |
|   News, etc.)   |       |                 |       |                 |
+-----------------+       +--------^--------+       +-----------------+
                                   |
                                   | (SQLite DB)
                                   |
                               +---+---+
                               |  Cache  |
                               | (news.db) |
                               +---------+
```

## 2. Backend (Go Lang)

The backend is a lightweight Go application responsible for:

*   **News Aggregation:** Periodically fetching articles from configured RSS feeds.
*   **Data Caching:** Storing fetched articles in a local SQLite database (`news.db`).
*   **API Serving:** Exposing RESTful endpoints for the frontend to consume.

### Key Components:

*   **`main.go`**: The entry point of the application. It sets up the HTTP server, defines routes, and initializes the background news caching job. It dynamically binds to a port specified by the `PORT` environment variable (defaulting to 8080 for local development).
*   **`db/db.go`**: Handles all database interactions. It initializes the SQLite database, defines functions for inserting and retrieving articles, and contains the `StartCachingJob` function which periodically fetches and caches news. It also includes the keyword-based ranking logic.
*   **`handlers/handlers.go`**: Contains the HTTP handler functions for the API endpoints, such as `/news` (to retrieve articles) and `/ad` (for advertisement data). It also includes a `/healthz` endpoint for health checks.
*   **`models/models.go`**: Defines the data structures (e.g., `NewsArticle`, `Ad`) used throughout the application.

### Features Implemented:

*   **RSS Feed Parsing:** Fetches and parses articles from a predefined list of RSS URLs.
*   **Local Caching:** Stores articles in `news.db` to reduce external requests and provide faster access.
*   **Keyword-Based Ranking:** Articles are assigned a rank based on the presence of specific cybersecurity keywords in their title and description.
*   **Filtering & Sorting:** The `/news` endpoint supports filtering by source, limiting the number of articles, date range, and sorting by publication date or calculated rank.
*   **Health Check Endpoint (`/healthz`):** A simple endpoint returning "OK" for monitoring and keep-alive services.

### Deployment on Render:

The backend is designed for deployment on platforms like Render's free tier. Key considerations include:

*   **Dynamic Port Binding:** Listens on the `PORT` environment variable provided by the hosting platform.
*   **Ephemeral Storage:** The `news.db` file is stored locally on the server's filesystem, which is ephemeral on Render's free tier. This means the database will reset on every server restart (e.g., due to inactivity or new deployments). The background caching job will repopulate it.

## 3. Frontend (Web / Android)

### Web Frontend (`news-api/test/index.html`):

*   A simple HTML page with embedded JavaScript and CSS.
*   Fetches news and ad data from the backend API.
*   Provides controls for source selection, article limit, date range, and sorting (including by rank).
*   Displays articles in a card-like format.

### Android Application (Planned):

The Android application, currently in its early stages, will serve as the primary user interface. It will interact with the same backend API.

## 4. Data Flow

1.  **Backend Startup:** On startup, the backend initializes its SQLite database and immediately starts a background job to fetch and cache news from all configured RSS sources. This job runs every 15 minutes.
2.  **Frontend Request:** When a user opens the web interface or the Android app, it sends an HTTP request to the backend's `/news` endpoint (e.g., `http://your-render-url.onrender.com/news?sortBy=rank`).
3.  **Backend Response:** The backend queries its local `news.db` based on the request parameters (filters, limits, sorting) and returns the relevant articles as a JSON payload.
4.  **Frontend Display:** The frontend receives the JSON data and renders the articles for the user.

## 5. Future Vision: The Cybersecurity Knowledge Hub

Beyond simple news aggregation, the long-term vision is to evolve the application into a comprehensive, actionable **Cybersecurity Knowledge Hub**. This involves transforming the aggregated data into a structured, interconnected resource that provides not just news, but also context, tutorials, and vulnerability intelligence.

**Confidence Score: 8/10** - This is a highly achievable vision with a phased approach.

### Core Components of the Hub:

1.  **Multiple Content Categories:** Expand beyond news to include:
    *   **Tech News:** General technology news to broaden the app's appeal.
    *   **Tutorials & How-Tos:** Practical, hands-on guides for security tasks and learning.
    *   **Vulnerability Disclosures:** A dedicated feed for the latest CVEs and security advisories.
2.  **Intelligent Reporting:** Generate a "Daily Threat Report" that synthesizes the day's key events, analyzes the overall threat level, and provides actionable advice.
3.  **Platform-Specific Filtering:** Allow users to filter all content by platform (e.g., "Linux", "Windows", "Android") to get a tailored view of relevant news, tutorials, and vulnerabilities.

### Phased Action Plan:

*   **Phase 1 (Short-Term): Add Content Categories.**
    *   **Action:** Implement the "Tech News" category as a proof-of-concept.
    *   **Backend:** Add a `category` field to the database and API. Find and integrate high-quality RSS feeds for tech news.
    *   **Frontend:** Add a UI element (e.g., tabs, dropdown) to switch between categories.

*   **Phase 2 (Mid-Term): Introduce Vulnerability & Tutorial Feeds.**
    *   **Action:** Research and integrate reliable sources for CVEs (e.g., NIST NVD) and security tutorials.
    *   **Backend:** Enhance the data model to differentiate between article types (news, vulnerability, tutorial). Implement more sophisticated tagging and categorization.

*   **Phase 3 (Long-Term): Implement Intelligent Cross-Linking & Reporting.**
    *   **Action:** Develop a system to connect related items (e.g., link a new vulnerability to a tutorial on how to fix it).
    *   **Backend:** This may require more advanced NLP/AI models to identify relationships between articles. Implement the "Daily Threat Report" generation logic.
    *   **Frontend:** Design a dedicated UI to present the daily report and showcase linked content.

## 6. AI-Powered Summarization (Cloud Function)

To provide significant user value without overloading the main backend, on-demand article summarization will be implemented using a serverless cloud function.

**Confidence Score: 9/10** - This is a standard, robust pattern for offloading intensive, spiky workloads.

### Architecture Flow:

```
+-------------+   1. Request   +------------------+   2. Invoke   +--------------------+
| Android App |-------------->|  Go Backend API  |------------->| Go Cloud Function  |
| (User taps  |   Summary    | (e.g., /summarize) |             | (Serverless)       |
| "Summarize")|              +------------------+             +---------+----------+
+-------------+                        |                                  | 3. Scrape & Call AI
      ^                              | 5. Return                          |
      | 6. Display                   |    Summary                         |
      +------------------------------+                                   v
                                                                +-----------------+
                                                                |   Gemini API    |
                                                                +-----------------+
```

### Key Implementation Details:

*   **On-Demand Trigger:** The process is initiated by a user action (e.g., tapping a button), not run automatically. This is crucial for managing costs and API limits.
*   **Go-Based Cloud Function:** The serverless function will be written in Go for maximum performance and efficiency.
*   **API Key Security:** The Gemini API key will be stored as a secure environment variable within the cloud function's configuration, never in the app or backend repository.
*   **Cost Control:**
    *   **Rate Limiting:** The Go backend will enforce a strict rate limit on the `/summarize` endpoint to prevent abuse.
    *   **Caching:** The cloud function will cache summaries for a set period (e.g., 24 hours) to avoid re-processing the same URL, significantly reducing API calls.

## 7. Database Migration Strategy

To prepare for a future move to a persistent database and to follow best practices, a simple database migration system will be implemented in the Go backend.

*   **Method:** On startup, the application will check a `schema_migrations` table and apply any new, versioned `.sql` script files found in a `migrations/` directory.
*   **Current State:** While not strictly necessary with Render's ephemeral filesystem, this makes the application more robust and production-ready.

## 8. Android Client-Side Features (Future)

The following features will be implemented directly within the Android application, leveraging its local storage capabilities.

*   **Local Database & Persistence:** The app will use a local SQLite database (via Room) to store articles, enabling offline access and a persistent user-specific archive.
*   **Export to CSV:** A feature will be added to allow users to export the articles from their local database into a CSV file for external use.
*   **Personal Notes:** Users will be able to create and attach personal notes to specific articles. This will require a new table in the local app database to store the notes and their relationship to articles.

## 9. Development Strategy & Order of Operations

A **backend-first** approach will be adopted for implementing new features. This ensures that the API contract is clearly defined before any frontend work begins, preventing rework.

### Data Synchronization Strategy:

Features requiring data to be shared between clients (e.g., the Android app and the web interface) such as **Personal Notes** would require user authentication (e.g., OAuth) and a persistent, centralized database. 

**Decision:** The implementation of a full authentication system is a significant undertaking and is deferred for a later stage. For the foreseeable future, features like "Personal Notes" and "CSV Export" will be implemented as **device-specific functionalities**. Data created on a device will remain on that device.

## 10. UI/UX Polish (Upcoming Tasks)

*   **Web Frontend: Fix Date Picker Visibility:** The date input fields in the `index.html` file are currently transparent, making the selected date invisible to the user. This needs to be fixed to provide clear visual feedback.

## 11. Other Future Enhancements

*   **Client-Side Persistence (Android):** Implement a local SQLite database on the Android device to store fetched articles. This would enable:
    *   **Offline Access:** Users can view previously fetched articles without an internet connection.
    *   **User-Specific Archive:** A persistent history of articles on the user's device, independent of server restarts.
    *   **Export Functionality:** Allow users to export their local archive to a CSV file.
*   **AI-Powered Summarization:**
    *   **Challenge:** RSS feeds often provide only snippets.
    *   **Solution:** Implement a web scraping component in the backend to fetch the full article content from source URLs (respecting `robots.txt` and terms of service).
    *   **Model Integration:** Use a small, efficient language model (e.g., Gemma 2B, Gemini Nano) to generate concise summaries from the full article text.
*   **Semantic Similarity Search:**
    *   **Concept:** Allow users to find articles semantically similar to one they are viewing.
    *   **Implementation:** Use a small, on-device sentence embedding model (e.g., a fine-tuned Sentence-BERT variant) to generate vector embeddings for article titles/descriptions. Store these embeddings in the client-side SQLite database. Perform cosine similarity calculations on the device.
*   **Keep-Alive Service:** Utilize an external cron job service (e.g., cron-job.org) to periodically ping the backend's `/healthz` endpoint. This prevents the Render free-tier service from going to sleep due to inactivity.
*   **Improved Error Handling & Logging:** Enhance logging for better debugging and monitoring in a production environment.
*   **User Authentication/Personalization:** For a more advanced application, implement user accounts to save preferences, custom feeds, or personalized recommendations.
*   **Monetization:** Integrate more sophisticated ad serving or premium features.
