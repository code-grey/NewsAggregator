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

## API Endpoints

*   **`/news`**: Fetches and returns news articles. Supports `source`, `limit`, `start`, and `end` query parameters.

## Environment Variables

- **`PORT`**: The port on which the server will listen. Defaults to `8080`.
- **`APP_URL`** (Optional but Recommended): The publicly accessible URL of your deployed application (e.g., `https://your-app.onrender.com`). If provided, the application will ping its own `/healthz` endpoint every 4 minutes to prevent it from sleeping on free hosting tiers.

## Security Considerations

This API includes basic security measures such as rate limiting and security headers. For production deployment, it is highly recommended to deploy this API behind a reverse proxy (e.g., Nginx, Caddy) to handle TLS encryption (HTTPS).
