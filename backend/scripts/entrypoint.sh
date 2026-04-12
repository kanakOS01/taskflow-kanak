#!/bin/sh

set -e

echo "Waiting for database to be ready..."
until psql "$DATABASE_URL" -c '\q' > /dev/null 2>&1; do
  echo "Database is unavailable - sleeping"
  sleep 1
done

echo "Running migrations..."
if command -v migrate > /dev/null 2>&1; then
    migrate -path /app/migrations -database "$DATABASE_URL" up
else
    echo "migrate tool not found. Ensure it is installed in the image."
    exit 1
fi

echo "Seeding database..."
if [ -f "/app/scripts/seed.sql" ]; then
    psql "$DATABASE_URL" -f /app/scripts/seed.sql
    echo "Seeding completed."
else
    echo "Seed script not found at /app/scripts/seed.sql, skipping."
fi

echo "Starting application..."
exec ./taskflow
