#!/bin/bash

# Script to run tests with filtered coverage
# Excludes: mocks, generated files (pb.go, sqlc-generated), view.go, main.go

set -e

# Run tests and generate coverage
echo "Running tests with coverage..."
go test -coverprofile=coverage.out ./...

# Check if coverage.out exists
if [ ! -f coverage.out ]; then
    echo "Error: coverage.out not found"
    exit 1
fi

# Filter coverage.out
# Excludes:
# - mocks/ - mock implementations
# - .pb.go - protobuf generated files
# - query.sql.go - sqlc generated queries
# - repository/db/models.go - sqlc generated models
# - repository/db/db.go - sqlc generated db interface
# - view.go - UI rendering code
# - main.go - entry points
echo "Filtering coverage..."
grep -v -E '(mocks/|\.pb\.go|query\.sql\.go|repository/db/models\.go|repository/db/db\.go|view\.go|main\.go)' coverage.out > coverage.filtered.out

# Show filtered coverage
echo ""
echo "=== Filtered Coverage Report ==="
echo ""
go tool cover -func=coverage.filtered.out

# Calculate and show total coverage
echo ""
echo "=== Total Coverage (Filtered) ==="
total_coverage=$(go tool cover -func=coverage.filtered.out | grep total | awk '{print $3}')
echo "Total: $total_coverage"

# Optional: Generate HTML report
if [ "$1" = "--html" ]; then
    echo ""
    echo "Generating HTML report..."
    go tool cover -html=coverage.filtered.out -o coverage.html
    echo "HTML report generated: coverage.html"
fi