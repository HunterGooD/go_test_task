include .env
$(eval export $(shell sed -ne 's/ *#.*$$//; /./ s/=.*$$// p' .env))
MIGRATION_NAME ?= new_migration
dbu:
	goose -dir db/migrations up

dbd:
	goose -dir db/migrations down

dbc:
	goose -dir db/migrations create "$(MIGRATION_NAME)" sql

compose_bd:
	docker-compose up
