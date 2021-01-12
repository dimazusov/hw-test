build:
	go build -o ./bin/calendar ./cmd/calendar/main.go

run:
	./bin/calendar

test:
	go test -race ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.30.0

lint: install-lint-deps
	golangci-lint run ./...

migrate:
	cd migrations && goose postgres "user=db password=db dbname=db sslmode=disable" up

.PHONY: build test lint
