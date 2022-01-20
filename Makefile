build:
	go build -v ./cmd/app

up:
	docker-compose -f ./docker/docker-compose.yml --env-file .env up -d

down:
	docker-compose -f ./docker/docker-compose.yml down

.DEFAULT_GOAL := build