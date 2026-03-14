#!/bin/sh
set -e

echo "Running database migrations..."
migrate -path /app/migrations -database "mysql://root:root@tcp(db:3306)/expense_tracker" up

echo "Starting server..."
exec ./server
