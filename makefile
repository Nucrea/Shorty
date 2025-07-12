.PHONY: run-dev
run-deps:
	./deploy/run-deps.sh

run: run-deps
	go run ./cmd/shorty --env ./deploy/dev.env

.PHONY: pprof
pprof:
	curl -XPOST http://localhost:8081/profile/start -H authorization:testapikey

.PHONY: pprof-stop
pprof-stop:
	curl -XPOST http://localhost:8081/profile/stop -H authorization:testapikey

.PHONY: pgclear
pgclear:
	docker compose down -v postgres && docker compose up -d postgres
