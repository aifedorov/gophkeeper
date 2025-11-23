.PHONY: build test run_server docker-up docker-down lint fmt all

fmt:
	@echo "Formatting..."
	go fmt ./...

build:
	@echo "Building..."
	(cd cmd/gophkeeper && go build -buildvcs=false -o gophkeeper)
	@echo "Build complete: cmd/gophkeeper/gophkeeper"

test:
	@echo "Running unit tests..."
	go test ./...

run_server:
	(cd cmd/gophkeeper && go run main.go)

docker-up:
	docker-compose up --build

docker-down:
	docker-compose down

lint:
	@echo "Running staticlint..."
	go run ./cmd/staticlint/main.go ./...

all: fmt lint build test
	@echo "All checks completed successfully!"
