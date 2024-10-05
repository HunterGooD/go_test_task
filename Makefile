include .env
$(eval export $(shell sed -ne 's/ *#.*$$//; /./ s/=.*$$// p' .env))

APP ?= music_service
BUILD ?= dev
MIGRATION_NAME ?= new_migration
OUTPUT_DIR ?= ./dist/${BUILD}/bin
.PHONY: cleango buildgo

dbu:
	goose -dir db/migrations up

dbd:
	goose -dir db/migrations down

dbc:
	goose -dir db/migrations create "$(MIGRATION_NAME)" sql

compose_db: 
	$(eval export HOST_DB=postgres)
	export DB_CONNECTION="postgresql://${USERNAME}:${PASSWORD}@${HOST_DB}:5432/${DB_NAME}?sslmode=disable" && \
	export GOOSE_DBSTRING=${DB_CONNECTION} && \
	docker-compose up

compose_db_sh:
	$(eval export HOST_DB=postgres)
	export DB_CONNECTION="postgresql://${USERNAME}:${PASSWORD}@${HOST_DB}:5432/${DB_NAME}?sslmode=disable" && \
	export GOOSE_DBSTRING=${DB_CONNECTION} && \
	docker-compose run --rm goose /bin/sh 
	# /go/bin/goose -dir=/app/db/migrations/ -v status

cleango:
	@rm -rf ${OUTPUT_DIR}

buildgo: cleango
	@mkdir -p ${OUTPUT_DIR}
	CGO_ENABLED=0 GOOS=linux go build -o ${OUTPUT_DIR}/${APP} ./cmd/music_service/

rungo:
	go run ./cmd/music_service/ 

testgo:
	go test ./...
