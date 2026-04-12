# Build Stage
FROM golang:1.26-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o taskflow cmd/api/main.go

# Runtime Stage
FROM alpine:3.18

# Installing migrate and postgresql-client
RUN apk add --no-cache curl postgresql-client \
    && curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xvz \
    && mv migrate /usr/local/bin/

WORKDIR /app
COPY --from=builder /app/taskflow .
COPY migrations /app/migrations
COPY scripts/seed.sql /app/scripts/seed.sql
COPY scripts/entrypoint.sh /app/entrypoint.sh

RUN chmod +x /app/entrypoint.sh

EXPOSE 8000

ENTRYPOINT ["/app/entrypoint.sh"]
