.PHONY: server client build-server build-client build-client-all docker-up docker-down docker-restart docker-logs docker-db-up docker-clean test test-cover lint fmt all env

# Development
server:
	DATABASE_URI="postgres://gophkeeper:password@localhost:5432/gophkeeper?sslmode=disable" \
	GRPC_ADDRESS="localhost:50051" \
	LOG_LEVEL="debug" \
	go run ./cmd/server/main.go

client:
	SERVER_ADDRESS="localhost:50051" \
	go run ./cmd/client/main.go

# Docker
docker-up:
	docker-compose up --build

docker-down:
	docker-compose down

docker-restart:
	docker-compose restart server

docker-logs:
	docker-compose logs -f server

docker-db-up:
	docker-compose up -d postgres migrate

docker-clean:
	docker-compose down -v
	docker volume prune -f

# Build
build-server:
	@mkdir -p dist
	cd cmd/server && go build -buildvcs=false -o ../../dist/gophkeeper-server main.go

build-client:
	@mkdir -p dist
	cd cmd/client && go build -buildvcs=false -o ../../dist/gophkeeper-client main.go

build-client-all:
	@mkdir -p dist
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -buildvcs=false -ldflags="-s -w" -o dist/gophkeeper-client-linux-arm64 ./cmd/client/main.go
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -buildvcs=false -ldflags="-s -w" -o dist/gophkeeper-client-darwin-amd64 ./cmd/client/main.go
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -buildvcs=false -ldflags="-s -w" -o dist/gophkeeper-client-windows-amd64.exe ./cmd/client/main.go

# Testing
test:
	go test -v ./...

test-cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

lint:
	golangci-lint run ./...

fmt:
	go fmt ./...

# Run all checks
all: fmt lint test
