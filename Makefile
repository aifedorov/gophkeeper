.PHONY: build test run_server docker-db-up docker-db-down docker-up docker-down lint fmt all

fmt:
	go fmt ./...

build:
	(cd cmd/server && go build -buildvcs=false -o gophkeeper)

test:
	go test ./...

run_server:
	(cd cmd/gophkeeper && go run main.go)

docker-db-up:
	docker-compose up -d postgres migrate

docker-db-down:
	docker-compose stop postgres migrate

docker-up:
	docker-compose up --build

docker-down:
	docker-compose down

lint:
	go run ./cmd/staticlint/main.go ./...

all: fmt lint build test
	@echo "All checks completed successfully!"
