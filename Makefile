.PHONY: run build test cover tidy docker-build up down logs
IMAGE ?= rest-api:dev
.DEFAULT_GOAL := help

run:    
	go run .
build: 
	go build -o bin/server .
test:  
	go test -count=1 ./...
cover: 
	go test -count=1 -coverprofile=cover.out ./...
	go tool cover -func=cover.out | tail -1
tidy: 
	go mod tidy && go mod verify
docker-build: 
	docker build -t $(IMAGE) .
up:   
	docker compose up --build -d
down: 
	docker compose down
logs: 
	docker compose logs -f api