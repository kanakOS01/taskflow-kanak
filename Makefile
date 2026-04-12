.PHONY: build run test migrate-up migrate-down db api seed help

DATABASE_URL ?= postgres://test:test@localhost:5432/taskflow?sslmode=disable

help:
	@echo "Available commands:"
	@echo "  make build          : Build the API server"
	@echo "  make run            : Run the API server directly"
	@echo "  make migrate-up     : Run all DB migrations UP"
	@echo "  make migrate-down   : Run one DB migration DOWN"
	@echo "  make seed           : Seed the DB with required test data"

build:
	go build -o bin/api cmd/api/main.go

run:
	go run cmd/api/main.go

migrate-up:
	migrate -path migrations -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path migrations -database "$(DATABASE_URL)" down 1

seed:
	psql "$(DATABASE_URL)" -f scripts/seed.sql
