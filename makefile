.PHONY: run-dev
run-dev:
	docker compose stop alloy || true
	yes | docker compose rm alloy -v || true

	rm -rf .run
	mkdir -p .run
	touch .run/shorty.log
	chmod 777 .run/shorty.log

	docker compose up -d
	
	SHORTY_MINIO_ENDPOINT=localhost:9000 \
	SHORTY_MINIO_ACCESS_KEY=miniokey \
	SHORTY_MINIO_ACCESS_SECRET=miniokey \
	SHORTY_LOG_FILE=".run/shorty.log" \
	SHORTY_APP_PORT=8081 \
	SHORTY_APP_URL=http://localhost:8081 \
	SHORTY_OPENTELEMETRY_URL="http://localhost:4318" \
	SHORTY_POSTGRES_URL=postgres://postgres:postgres@localhost:5432/postgres \
	SHORTY_REDIS_URL=redis://localhost:6379 \
	go run .

.PHONY: pgclear
pgclear:
	docker compose down -v postgres && docker compose up -d postgres
