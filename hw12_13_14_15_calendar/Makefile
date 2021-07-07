BIN := "./bin/calendar"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar

run: build
	$(BIN) --config ./configs/config.toml

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	$(BIN) version

test:
	go test -v -race ./...

up:
	goose -dir migrations postgres "user=calendar password=calendar dbname=calendar sslmode=disable" up

down:
	goose -dir migrations postgres "user=calendar password=calendar dbname=calendar sslmode=disable" down

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.40.1

lint: install-lint-deps
	golangci-lint run ./...

generate:
	go generate ./...

fix:
	golangci-lint run --fix

.PHONY: build run build-img run-img version test lint up down generate fix