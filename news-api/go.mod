module news-api

go 1.24.4

replace (
	"news-api/handlers" => "./handlers"
	"news-api/models" => "./models"
	"news-api/db" => "./db"
)

require (
	github.com/PuerkitoBio/goquery v1.8.0 // indirect
	github.com/andybalholm/cascadia v1.3.1 // indirect
	github.com/aymerick/douceur v0.2.0 // indirect
	github.com/gorilla/css v1.0.1 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/mattn/go-sqlite3 v1.14.30 // indirect
	github.com/microcosm-cc/bluemonday v1.0.27 // indirect
	github.com/mmcdole/gofeed v1.3.0 // indirect
	github.com/mmcdole/goxpp v1.1.1-0.20240225020742-a0c311522b23 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	golang.org/x/net v0.26.0 // indirect
	golang.org/x/text v0.16.0 // indirect
	golang.org/x/time v0.12.0 // indirect
)
