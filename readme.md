<!-- for build api -->
go build -o bin/restaurant ./cmd/api

<!-- for cross platform build [run at gitbash]-->
GOOS=linux GOARCH=amd64 go build -o restaurant ./cmd/api

<!-- for run api with prod -->
APP_ENV=dev ./bin/restaurant

<!-- for run without making build -->
APP_ENV=dev go run ./cmd/api

<!-- droplet restart go server service linux cmd-->
systemctl restart restaurant.service
systemctl status restaurant.service