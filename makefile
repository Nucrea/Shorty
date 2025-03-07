.PHONY: run-dev
run-dev:
	SHORTY_APP_PORT=8081 \
	SHORTY_APP_URL=http://localhost:8081 \
	SHORTY_POSTGRES_URL=postgres://postgres:postgres@localhost:5432/postgres \
	SHORTY_REDIS_URL=redis://localhost:6379 \
	go run .
