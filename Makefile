DB_URL="postgres://dev:123456789@localhost:5433/facebookad?sslmode=disable"

.PHONY: start-services migration-new migration-up migration-down

start-services:
	docker-compose -f ./scripts/docker-compose.yml up -d

migration-new:
	dbmate --migrations-dir=./migrations new $(name)
migration-up:
	dbmate --migrations-dir=./migrations --url=$(DB_URL) up
migration-down:
	dbmate --migrations-dir=./migrations --url=$(DB_URL) down

