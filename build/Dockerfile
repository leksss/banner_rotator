FROM golang:1.16.2 as build

ENV CODE_DIR /go/src
ENV BIN_FILE ${CODE_DIR}/bin/banner_rotator

WORKDIR ${CODE_DIR}

# Кэшируем слои с модулями
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . ${CODE_DIR}

# Собираем статический бинарник Go (без зависимостей на Си API), иначе он не будет работать в alpine образе.
ARG LDFLAGS
RUN CGO_ENABLED=0 go build \
        -ldflags "$LDFLAGS" \
        -o ${BIN_FILE} cmd/banner_rotator/*

FROM alpine:3.9

RUN apk add --no-cache openssl

ENV DOCKERIZE_VERSION v0.6.1
RUN wget https://github.com/jwilder/dockerize/releases/download/$DOCKERIZE_VERSION/dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && tar -C /usr/local/bin -xzvf dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && rm dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz

ENV CODE_DIR /go/src
ENV BIN_FILE ${CODE_DIR}/bin/banner_rotator
ENV CONFIG_FILE ${CODE_DIR}/configs/config_docker.yaml
ENV LOG_FILE ${CODE_DIR}/logs/logger.json
ENV START_FILE ${CODE_DIR}/build/start.sh

COPY --from=build ${BIN_FILE} ${BIN_FILE}
COPY --from=build ${CONFIG_FILE} ${CONFIG_FILE}
COPY --from=build ${LOG_FILE} ${LOG_FILE}
COPY --from=build ${START_FILE} ${START_FILE}

EXPOSE 8080 8081

CMD ${START_FILE}
