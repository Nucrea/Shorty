.PHONY: run-dev
run-dev:
	SHORTY_POSTGRES_URL=postgres://postgres:postgres@localhost:5432/postgres \
	SHORTY_REDIS_URL=redis://localhost:6379 \
	SHORTY_BASE_URL=http://localhost:8081 \
	go run .
