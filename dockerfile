FROM golang:1.23-alpine
WORKDIR /app

COPY . .

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

ENV SHORTY_POSTGRES_URL=""
ENV SHORTY_BASE_URL=""

EXPOSE 8081

CMD ["./app"]