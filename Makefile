.PHONY: build test run_server docker-db-up docker-db-down docker-migrate-up docker-migrate-down docker-up docker-down lint fmt all

fmt:
	go fmt ./...

build:
	(cd cmd/server && go build -buildvcs=false -o gophkeeper)

test:
	go test ./...

run_server:
	(cd cmd/gophkeeper && go run main.go)

docker-db-up:
	docker-compose up -d postgres

docker-db-down:
	docker-compose stop postgres

docker-migrate-up:
	migrate -path ./migrations -database 'postgres://gophkeeper:password@localhost:5432/gophkeeper?sslmode=disable' up

docker-migrate-down:
	migrate -path ./migrations -database 'postgres://gophkeeper:password@localhost:5432/gophkeeper?sslmode=disable' down

docker-up:
	docker-compose up --build

docker-down:
	docker-compose down

lint:
	go run ./cmd/staticlint/main.go ./...

all: fmt lint build test
	@echo "All checks completed successfully!"
