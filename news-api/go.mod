module news-api

go 1.24.4

replace (
	news-api/db => ./db
	news-api/handlers => ./handlers
	news-api/models => ./models
)

require (
	github.com/mattn/go-sqlite3 v1.14.30
	github.com/microcosm-cc/bluemonday v1.0.27
	github.com/mmcdole/gofeed v1.3.0
	github.com/pemistahl/lingua-go v1.4.0
	github.com/stretchr/testify v1.11.1
	golang.org/x/time v0.12.0
)

require (
	github.com/PuerkitoBio/goquery v1.8.0 // indirect
	github.com/andybalholm/cascadia v1.3.1 // indirect
	github.com/aymerick/douceur v0.2.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/gorilla/css v1.0.1 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/mmcdole/goxpp v1.1.1-0.20240225020742-a0c311522b23 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/shopspring/decimal v1.3.1 // indirect
	golang.org/x/exp v0.0.0-20221106115401-f9659909a136 // indirect
	golang.org/x/net v0.26.0 // indirect
	golang.org/x/text v0.16.0 // indirect
	google.golang.org/protobuf v1.31.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
