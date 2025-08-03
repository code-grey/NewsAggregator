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

## 5. Future Enhancements

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
