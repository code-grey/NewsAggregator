# ThreatFeed Backend API

This is the Go-based backend API for the ThreatFeed news aggregator. It fetches and caches cybersecurity news articles from various RSS feeds into a local SQLite database and serves them via a REST API.

## Removed Sources

**Dark Reading (`https://www.darkreading.com/rss/all.xml`)**

This source has been removed from the list of active RSS feeds due to persistent `403 Forbidden` errors when attempting to fetch its content. This indicates that the server is actively blocking automated requests, likely as an anti-scraping measure. While more advanced techniques (e.g., using proxies, headless browsers) could potentially bypass this, they are beyond the scope and complexity desired for this simple prototype.

## Getting Started

1.  **Ensure Go is installed:** [https://golang.org/doc/install](https://golang.org/doc/install)
2.  **Navigate to this directory:** `cd news-api`
3.  **Run for development:** `go run main.go`

    The API will start on `http://localhost:8080`.

    *Note: The first run will populate the `news.db` SQLite file, which might take a few moments as it fetches articles from all sources.*

## Building for Production

To create a smaller, optimized binary for production, use the following build command. This strips debug information and reduces the file size significantly.

```bash
go build -ldflags="-s -w" -o news-api-prod main.go
```

This will create a `news-api-prod` executable in the current directory.

## API Documentation

This API provides endpoints to retrieve news articles and a daily threat assessment. All endpoints return responses in JSON format.

### Get News Articles

- **Endpoint:** `/news`
- **Method:** `GET`
- **Description:** Fetches a list of news articles. Articles can be filtered by source, category, a search query, and date range.

#### Query Parameters

| Parameter | Type    | Description                                                                                                  | Example                               |
| :-------- | :------ | :----------------------------------------------------------------------------------------------------------- | :------------------------------------ |
| `source`  | string  | Filter articles by a specific RSS feed URL.                                                                  | `?source=https://www.bleepingcomputer.com/feed/` |
| `category`| string  | Filter articles by category. Supported values are `Cybersecurity`, `Tech`, and `Defense`.                      | `?category=Cybersecurity`             |
| `search`  | string  | A search term to filter articles by title or description. The search is case-insensitive.                      | `?search=ransomware`                  |
| `limit`   | integer | The maximum number of articles to return. Defaults to `20`.                                                    | `?limit=10`                           |
| `start`   | string  | The start date for filtering articles, in `YYYY-MM-DD` format.                                               | `?start=2023-10-26`                   |
| `end`     | string  | The end date for filtering articles, in `YYYY-MM-DD` format.                                                 | `?end=2023-10-27`                     |
| `sortBy`  | string  | The sorting order for the articles. Supported values are `publishedAt` (default) and `rank`.                 | `?sortBy=rank`                        |

#### Example Request (Using `curl`)

```bash
curl "http://localhost:8080/news?category=Cybersecurity&search=vulnerability&limit=5"
```

#### Example Response

```json
[
    {
        "id": 123,
        "title": "Critical Vulnerability Found in Popular Web Server",
        "description": "A critical vulnerability has been discovered...",
        "imageUrl": "https://example.com/image.png",
        "url": "https://example.com/article",
        "sourceUrl": "https://feeds.feedburner.com/TheHackersNews",
        "publishedAt": "2023-10-27T10:00:00Z",
        "rank": 5,
        "category": "Cybersecurity"
    }
]
```

### Get Today's Threat Score

- **Endpoint:** `/today-threat`
- **Method:** `GET`
- **Description:** Provides a threat assessment based on the articles published in the last 24 hours. The threat level is categorized as `Code Red`, `Attention`, or `Business as Usual`.

#### Example Request (Using `curl`)

```bash
curl "http://localhost:8080/today-threat"
```

#### Example Response

```json
{
    "lowRankCount": 15,
    "mediumRankCount": 5,
    "highRankCount": 2,
    "totalArticles": 22,
    "threatLevel": "Code Red"
}
```

## Environment Variables

- **`PORT`**: The port on which the server will listen. Defaults to `8080`.
- **`APP_URL`** (Optional but Recommended): The publicly accessible URL of your deployed application (e.g., `https://your-app.onrender.com`). If provided, the application will ping its own `/healthz` endpoint every 4 minutes to prevent it from sleeping on free hosting tiers.

## Security Considerations

This API includes basic security measures such as rate limiting and security headers. For production deployment, it is highly recommended to deploy this API behind a reverse proxy (e.g., Nginx, Caddy) to handle TLS encryption (HTTPS).
