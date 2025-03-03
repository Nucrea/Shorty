FROM golang:1.23-alpine
WORKDIR /app

COPY . .

RUN go install && go mod verify
RUN go build -o app
RUN chmod +x app

ENV SHORTY_POSTGRES_URL=""
ENV SHORTY_BASE_URL=""

EXPOSE 8081

CMD ["./app"]