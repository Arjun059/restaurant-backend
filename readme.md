<!-- for build api -->
go build -o bin/restaurant ./cmd/api

<!-- for run api with prod -->
APP_ENV=dev ./bin/restaurant

<!-- for run without making build -->
APP_ENV=dev go run ./cmd/api