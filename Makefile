SHELL := /bin/bash
GO ?= go
HTTP_ADDR ?= :8081
ACCESS_TOKEN_TTL ?= 15m
REFRESH_TOKEN_TTL ?= 168h
REFRESH_COOKIE_NAME ?= refresh_token
COOKIE_SECURE ?= true
export DATABASE_URL TEST_DATABASE_URL JWT_SECRET ACCESS_TOKEN_TTL REFRESH_TOKEN_TTL REFRESH_COOKIE_NAME COOKIE_SECURE HTTP_ADDR LOG_LEVEL LOG_FORMAT
.PHONY: help run build test vet race fmt fmt-check deps starter-check ci check test-acceptance test-db test-db-race up down status db-up db-down clean
help:
	@echo "Starter: make starter-check / make ci / make run"
	@echo "Задание: make test-acceptance; БД после реализации: make test-db test-db-race"
	@echo "Goose: make up / make down / make status; 002 пока содержит TODO"
run:
	$(GO) run ./cmd/api
build:
	@mkdir -p bin
	$(GO) build -buildvcs=false -o bin/sessions-api ./cmd/api
	$(GO) build -buildvcs=false -o bin/sessions-migrate ./cmd/migrate
test:
	$(GO) test ./...
vet:
	$(GO) vet ./...
race:
	$(GO) test -race ./...
fmt:
	gofmt -w $$(find cmd internal -type f -name '*.go')
fmt-check:
	@files="$$(gofmt -l $$(find cmd internal -type f -name '*.go'))"; if [ -n "$$files" ]; then echo "$$files"; exit 1; fi
deps:
	$(GO) list -deps ./... >/dev/null
starter-check: fmt-check vet
	$(GO) test -run '^$$' -tags=acceptance ./...
	$(GO) test -run '^$$' -tags=integration ./internal/repository/postgres
ci: starter-check deps test race build
check: ci
test-acceptance:
	$(GO) test -tags=acceptance -count=1 ./...
test-db:
	@test -n "$$TEST_DATABASE_URL" || { echo "Задайте отдельную TEST_DATABASE_URL"; exit 1; }
	$(GO) test -tags=integration -count=1 ./internal/repository/postgres
test-db-race:
	@test -n "$$TEST_DATABASE_URL" || { echo "Задайте TEST_DATABASE_URL"; exit 1; }
	$(GO) test -race -tags=integration -count=1 ./internal/repository/postgres
up:
	$(GO) run ./cmd/migrate up
# Один вызов откатывает только последнюю применённую миграцию.
down:
	$(GO) run ./cmd/migrate down
status:
	$(GO) run ./cmd/migrate status
db-up: up
db-down: down
clean:
	rm -f bin/sessions-api bin/sessions-migrate
