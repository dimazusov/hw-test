BIN := "./bin/calendar"
BIN_SHEDULER := "./bin/calendar_scheduler"
BIN_SENDER := "./bin/calendar_sender"

DOCKER_IMG="calendar:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar && \
	go build -v -o $(BIN_SHEDULER) -ldflags "$(LDFLAGS)" ./cmd/calendar_scheduler && \
	go build -v -o $(BIN_SENDER) -ldflags "$(LDFLAGS)" ./cmd/calendar_sender

run: build
	$(BIN) --config ./configs/config.yaml & \
	$(BIN_SHEDULER) -config ./configs/scheduler_config.yaml & \
	$(BIN_SENDER) -config ./configs/sender_config.yaml

run-with-docker:
	docker-compose -f deployments/docker-compose.yaml up -d

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
	go test -race ./internal/...

generate:
	protoc --proto_path=internal/server/grpc --go_out=internal/server/grpc/pb --go-grpc_out=internal/server/grpc/pb internal/server/grpc/*.proto

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.33.0

lint: install-lint-deps
	golangci-lint -c=./configs/golangci.yaml run ./...

migrate:
	cd migrations && goose postgres "user=db password=db dbname=db sslmode=disable" up

.PHONY: build test lint
.PHONY: build run build-img run-img version test lint
