BIN := "./bin/calendar"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar

build-sched:
	go build -v -o ./bin/scheduler ./cmd/scheduler

build-sender:
	go build -v -o ./bin/sender ./cmd/sender

run: build
	$(BIN) --config ./configs/calendar.toml

run-sched: build-sched
	./bin/scheduler --config ./configs/scheduler.toml

run-sender: build-sender
	./bin/sender --config ./configs/sender.toml

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
	goose -dir migrations postgres "user=calendar password=calendar dbname=postgres sslmode=disable" up

down:
	goose -dir migrations postgres "user=calendar password=calendar dbname=postgres sslmode=disable" down

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.40.1

lint: install-lint-deps
	golangci-lint run ./...

generate:
	go generate ./...

fix:
	golangci-lint run --fix

.PHONY: build run build-img run-img version test lint up down generate fix build-sched build-sender run-sched run-sender