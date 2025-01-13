start-services:
	docker-compose -f ./scripts/docker-compose.yml up -d

migration-new-%:
	dbmate --migrations-dir=./migrations new $*
migration-up:
	dbmate --migrations-dir=./migrations up
migration-down:
	dbmate --migrations-dir=./migrations rollback
