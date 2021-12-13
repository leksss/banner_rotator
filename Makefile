BIN := "./bin/banner_rotator"
MIGRATOR_BIN := "./bin/migrator"
DOCKER_IMG="banner_rotator:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	docker-compose -f build/docker-compose.yaml -f build/docker-compose.app.yaml build --no-cache

run:
	docker-compose -f build/docker-compose.yaml -f build/docker-compose.app.yaml up -d --build

run-external:
	docker-compose -f build/docker-compose.yaml up -d --build

run-local:
	go run ./cmd/banner_rotator/

stop:
	docker-compose -f build/docker-compose.yaml -f build/docker-compose.app.yaml down

stop-dev:
	docker-compose -f build/docker-compose.yaml down

ps:
	docker-compose -f build/docker-compose.yaml -f build/docker-compose.app.yaml ps

log:
	docker-compose -f build/docker-compose.yaml -f build/docker-compose.app.yaml logs -f

test:
	go test -race -count 100 ./internal/...

integration-test:
	go test ./... -v -race -tags=integration

version: build
	$(BIN) version

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.37.0

lint: install-lint-deps
	golangci-lint run ./...

build-migrator:
	go build -v -o $(MIGRATOR_BIN) -ldflags "$(LDFLAGS)" ./cmd/migration

migrate:
	$(MIGRATOR_BIN) -dir=./migrations mysql up

generate:
	go generate ./cmd/banner_rotator

.PHONY: build run build-img run-img version test lint

