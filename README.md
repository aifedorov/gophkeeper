# GophKeeper

A secure client-server password manager for storing credentials, bank cards, binary files, and text notes.

## Purpose

GophKeeper is a secure storage system that allows users to safely store and manage sensitive information:
- **Credentials** — login/password pairs
- **Bank Cards** — card number, expiry date, CVV, cardholder name
- **Binary Files** — any files with streaming upload/download
- **Text Notes** — arbitrary text data

All data supports custom metadata (website, bank name, notes, etc.) and is encrypted at rest.

## Tech Stack

| Component      | Technology                                               |
|----------------|----------------------------------------------------------|
| Language       | Go 1.25                                                  |
| Protocol       | gRPC with Protocol Buffers (binary protocol)             |
| Database       | PostgreSQL 17                                            |
| Migrations     | golang-migrate                                           |
| Authentication | JWT tokens                                               |
| Encryption     | AES-256-GCM (server-side encryption)                     |
| CLI Framework  | Cobra                                                    |
| Config         | Environment variables with `caarlos0/env`                |
| Logging        | Zap                                                      |
| Testing        | testify, go-mock                                         |
| Containerization| Docker, Docker Compose                                  |

## Architecture

```
├── cmd/
│   ├── client/          # CLI client entry point
│   └── server/          # gRPC server entry point
├── internal/
│   ├── client/
│   │   ├── application/ # App initialization
│   │   ├── cli/         # Cobra commands (auth, secrets)
│   │   ├── config/      # Client configuration
│   │   ├── domain/      # Business logic (auth, credential, binary, card, text)
│   │   └── infrastructure/
│   │       ├── grpc/    # gRPC client connections
│   │       └── storage/ # Local cache and session storage
│   └── server/
│       ├── api/grpc/    # Proto definitions and generated code
│       ├── application/ # App initialization
│       ├── config/      # Server configuration
│       ├── domain/      # Business logic
│       │   ├── auth/    # User registration, login, JWT
│       │   └── secret/  # Credential, Binary, Card services
│       └── infrastructure/
│           ├── crypto/  # AES encryption/decryption
│           ├── grpc/    # gRPC handlers and interceptors
│           ├── jwt/     # JWT token generation/validation
│           └── postgres/# Database connection
├── migrations/          # SQL migrations
└── pkg/                 # Shared packages (filestorage, logger, validator)
```

### Design Principles
- **Domain-Driven Design** — separation of domain, application, and infrastructure layers
- **Clean Architecture** — dependency inversion, interfaces for external dependencies
- **Repository Pattern** — database abstraction with sqlc-generated code

## Features

| Feature | Status |
|---------|--------|
| User registration and authentication | ✅ |
| JWT-based authorization | ✅ |
| Credentials (login/password) CRUD | ✅ |
| Bank cards CRUD | ✅ |
| Binary files with streaming upload/download | ✅ |
| Text notes CRUD | ✅ |
| Metadata support for all secret types | ✅ |
| Server-side AES-256 encryption | ✅ |
| Cross-platform CLI (Windows, Linux, macOS) | ✅ |
| Version and build info in client | ✅ |
| gRPC binary protocol | ✅ |
| Docker deployment | ✅ |
| Unit tests (80%+ coverage goal) | ✅ |

## Running the Server

### Prerequisites
- Docker and Docker Compose
- Go 1.25+ (for local development)

### Using Docker (Recommended)

```bash
# Start all services (PostgreSQL, migrations, server)
make docker-up

# Start only database with migrations
make docker-db-up

# View server logs
make docker-logs

# Stop all services
make docker-down
```

### Local Development

```bash
# Start database first
make docker-db-up

# Run server
make server
```

### Server Configuration

| Variable | Description | Default |
|----------|-------------|---------|
| `DATABASE_URI` | PostgreSQL connection string | required |
| `GRPC_ADDRESS` | gRPC listen address | required |
| `JWT_SECRET_KEY` | JWT signing secret | required |
| `JWT_EXPIRATION` | Token TTL | `24h` |
| `LOG_LEVEL` | Log level (debug/info/warn/error) | `info` |

## Using the Client

### Installation

Download pre-built binaries from `dist/` or build from source:

```bash
# Build for current platform
make build-client

# Build for all platforms (Linux, macOS, Windows)
make build-client-all
```

### Client Configuration

| Variable | Description | Default |
|----------|-------------|---------|
| `SERVER_ADDRESS` | gRPC server address | required |
| `LOG_LEVEL` | Log level | `info` |

### Commands

```bash
# Show version and build info
gophkeeper --version

# Show all available commands
gophkeeper commands

# Authentication
gophkeeper register -e <email> -p <password>
gophkeeper login -e <email> -p <password>

# Credentials
gophkeeper credential create -n <name> -l <login> -p <password> [-m <metadata>]
gophkeeper credential list
gophkeeper credential update -i <id> [-n <name>] [-l <login>] [-p <password>] [-m <metadata>]
gophkeeper credential delete -i <id>

# Bank Cards
gophkeeper card create -n <name> --number <number> --expiry <MM/YY> --cvv <cvv> --holder <name> [-m <metadata>]
gophkeeper card list
gophkeeper card update -i <id> [flags]
gophkeeper card delete -i <id>

# Binary Files
gophkeeper file upload -f <filepath> [-n <name>] [-m <metadata>]
gophkeeper file list
gophkeeper file download -i <id> -o <output_path>
gophkeeper file update -i <id> [-f <filepath>] [-n <name>] [-m <metadata>]
gophkeeper file delete -i <id>

# Text Notes
gophkeeper text create -n <name> -c <content> [-m <metadata>]
gophkeeper text list
gophkeeper text update -i <id> [-n <name>] [-c <content>] [-m <metadata>]
gophkeeper text delete -i <id>
```

### Example Session

```bash
# Set server address
export SERVER_ADDRESS="localhost:50051"

# Register new user
./gophkeeper register -e user@example.com -p mysecretpassword

# Login
./gophkeeper login -e user@example.com -p mysecretpassword

# Store a credential
./gophkeeper credential create -n "GitHub" -l "myuser" -p "mypass" -m "work account"

# List all credentials
./gophkeeper credential list

# Upload a file
./gophkeeper file upload -f ./secret.pdf -n "Tax Documents" -m "2024 taxes"

# Download a file
./gophkeeper file download -i <file-id> -o ./downloaded.pdf
```

## Tests and Documentation

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage report
make test-coverage
```

### Code Quality

```bash
# Run linter
make lint

# Format code
make fmt

# Run all checks (fmt, lint, test)
make all
```

### Generating Code

```bash
# Generate protobuf code
make proto-gen
```

### Test Coverage Target

The project aims for **80%+ test coverage** on domain business logic as per specification requirements.

### Documentation

All exported functions, types, and packages include GoDoc documentation. View documentation with:

```bash
go doc ./...
# or
godoc -http=:6060
```

## Build

```bash
# Build server binary
make build-server

# Build client binary (current platform)
make build-client

# Build client for all platforms
make build-client-all
# Outputs:
#   dist/gophkeeper-client-linux-arm64
#   dist/gophkeeper-client-darwin-amd64
#   dist/gophkeeper-client-windows-amd64.exe
```

## License

MIT
