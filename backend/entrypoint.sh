#!/bin/sh
set -e

echo "Running go mod tidy..."
go mod tidy

echo "Running database migrations..."
migrate -path /app/migrations -database "mysql://root:root@tcp(db:3306)/expense_tracker" up

echo "Starting server with Air (hot reload)..."
exec air
