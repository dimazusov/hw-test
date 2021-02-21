# Собираем в гошке
FROM golang:1.15.2 as build

ENV CODE_DIR /go/src/

WORKDIR ${CODE_DIR}

# Кэшируем слои с модулями
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . ${CODE_DIR}

RUN go test cmd/calendar_test/main_test.go
