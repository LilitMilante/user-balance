include ./.env

build:
	go build ./cmd/app

up:
	docker-compose -f ./docker/docker-compose.yml --env-file .env up -d

down:
	docker-compose -f ./docker/docker-compose.yml down


migrate-install:
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.15.1

migrate-up: migrate-install
	migrate -path ./migrations -database="${DB_SCHEME}://${DB_LOGIN}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable&&query" up

migrate-down: migrate-install
	migrate -path ./migrations -database="${DB_SCHEME}://${DB_LOGIN}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable&&query" down 1

migrate-new: migrate-install
	migrate create -ext sql -dir ./migrations "$(name)"


.DEFAULT_GOAL := build