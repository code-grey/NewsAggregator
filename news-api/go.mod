module news-api

go 1.24.4

replace (
	news-api/db => ./db
	news-api/handlers => ./handlers
	news-api/models => ./models
	news-api/gemini => ./gemini
)

require (
	github.com/mattn/go-sqlite3 v1.14.30
	github.com/microcosm-cc/bluemonday v1.0.27
	github.com/mmcdole/gofeed v1.3.0
	github.com/pemistahl/lingua-go v1.4.0
	golang.org/x/time v0.12.0
)

require (
	cloud.google.com/go v0.116.0 // indirect
	cloud.google.com/go/auth v0.9.3 // indirect
	cloud.google.com/go/compute/metadata v0.5.0 // indirect
	github.com/PuerkitoBio/goquery v1.8.0 // indirect
	github.com/andybalholm/cascadia v1.3.1 // indirect
	github.com/aymerick/douceur v0.2.0 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/google/s2a-go v0.1.8 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.3.4 // indirect
	github.com/gorilla/css v1.0.1 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/mmcdole/goxpp v1.1.1-0.20240225020742-a0c311522b23 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/shopspring/decimal v1.3.1 // indirect
	go.opencensus.io v0.24.0 // indirect
	golang.org/x/crypto v0.27.0 // indirect
	golang.org/x/exp v0.0.0-20221106115401-f9659909a136 // indirect
	golang.org/x/net v0.29.0 // indirect
	golang.org/x/sys v0.25.0 // indirect
	golang.org/x/text v0.18.0 // indirect
	google.golang.org/genai v1.28.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240903143218-8af14fe29dc1 // indirect
	google.golang.org/grpc v1.66.2 // indirect
	google.golang.org/protobuf v1.34.2 // indirect
)
