# Собираем в гошке
FROM golang:1.15.2 as build

ENV BIN_FILE /opt/calendar_test/calendar_test
ENV CODE_DIR /go/src/

WORKDIR ${CODE_DIR}

# Кэшируем слои с модулями
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . ${CODE_DIR}

# Собираем статический бинарник Go (без зависимостей на Си API),
# иначе он не будет работать в alpine образе.
ARG LDFLAGS
RUN CGO_ENABLED=0 go build \
        -ldflags "$LDFLAGS" \
        -o ${BIN_FILE} cmd/calendar_test/*

# На выходе тонкий образ
FROM alpine:3.9

LABEL ORGANIZATION="DMITRYI USOV INCORPRATED"
LABEL SERVICE="calendar_test"
LABEL MAINTAINERS="dimazusov@yandex.ru"

ENV BIN_FILE "/opt/calendar_test/calendar_test"
COPY --from=build ${BIN_FILE} ${BIN_FILE}

ENV CONFIG_FILE /etc/calendar/config_test.yaml
COPY ./configs/config.yaml ${CONFIG_FILE}

CMD ${BIN_FILE} -config ${CONFIG_FILE}
