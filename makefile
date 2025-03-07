.PHONY: run-dev
run-dev:
	docker compose stop alloy || true
	yes | docker compose rm alloy -v || true

	rm -rf .run
	mkdir -p .run
	touch .run/shorty.log
	chmod 777 .run/shorty.log

	docker compose up -d

	SHORTY_LOG_FILE=".run/shorty.log" \
	SHORTY_APP_PORT=8081 \
	SHORTY_APP_URL=http://localhost:8081 \
	SHORTY_ELASTICSEARCH_URL="http://localhost:9200" \
	SHORTY_POSTGRES_URL=postgres://postgres:postgres@localhost:5432/postgres \
	SHORTY_REDIS_URL=redis://localhost:6379 \
	go run .
