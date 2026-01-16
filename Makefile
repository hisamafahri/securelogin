.PHONY: dev format migrate-up migrate-down migrate-generate migrate-status migrate-hash

ifneq (,$(wildcard .env))
    include .env
    export
endif

DATABASE_URL ?= $(shell echo $$DATABASE_URL)
ifeq ($(DATABASE_URL),)
$(error DATABASE_URL environment variable is not set)
endif

dev:
	air

format:
	gofmt -w .
	golines --max-len=80 --base-formatter=gofumpt --shorten-comments -t 2 -w .
	golangci-lint run

migrate-up:
	atlas migrate apply --env gorm --url "$(DATABASE_URL)"

migrate-down:
	@read -p "How many migrations to rollback? (default: 1): " count; \
	count=$${count:-1}; \
	atlas migrate down --env gorm --url "$(DATABASE_URL)" $$count

migrate-generate:
	@read -p "Enter migration name: " name; \
	atlas migrate diff $$name --env gorm

migrate-status:
	atlas migrate status --env gorm --url "$(DATABASE_URL)"

migrate-hash:
	atlas migrate hash --env gorm
