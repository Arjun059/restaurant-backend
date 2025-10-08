# Define your application entry point
run:
	go run cmd/api/main.go

build:
	go build -o bin/api cmd/api/main.go

start: build
	./bin/api

test:
	go test ./tests -v

test.fresh:
	go clean -testcache
	go test ./tests

clean.build:
	rm -rf ./bin/**

clean.test:
	go clean -testcache

dev: 
	air

