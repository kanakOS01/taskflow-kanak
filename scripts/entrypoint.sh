#!/bin/sh

echo "Running migrations..."
# Check if migrate is available, otherwise assume production binary path
if command -v migrate > /dev/null 2>&1; then
    migrate -path /app/migrations -database "$DATABASE_URL" up
else
    echo "migrate tool not found. Ensure it is installed in the image."
    exit 1
fi

echo "Starting application..."
exec ./taskflow
