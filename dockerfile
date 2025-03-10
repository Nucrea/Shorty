FROM golang:1.23-alpine AS builder
WORKDIR /backend

COPY go.mod go.sum config.go main.go ./
COPY src ./src
COPY server ./server

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    go mod download -x

ENV GOCACHE=/root/.cache/go-build
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=cache,target="/root/.cache/go-build" \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    go build -ldflags "-s -w" -o app

RUN chmod +x app

FROM alpine:3.21 AS production
WORKDIR /backend

COPY --from=builder /backend/app app

ENV SHORTY_APP_PORT=8081
ENV SHORTY_APP_URL=""
ENV SHORTY_OPENTELEMETRY_URL=""
ENV SHORTY_POSTGRES_URL=""
ENV SHORTY_REDIS_URL=""

CMD ["./app"]