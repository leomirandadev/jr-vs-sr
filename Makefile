
.PHONY: docs

NAME = "capsulas"

include .env

docs:
	@swag init -g cmd/api/main.go --parseDependency

install: 
	@echo "installing goose..."
	@go install github.com/pressly/goose/v3/cmd/goose@v3.19.2
	@echo "installing swaggo..."
	@go install github.com/swaggo/swag/cmd/swag@v1.16.3
	@echo "installing air (hot reaload)..."
	@go install github.com/cosmtrek/air@v1.51.0
	@echo "downloading project dependencies..."
	@go mod tidy

build: 
	@echo $(NAME): Construindo a imagem
	@docker build --progress=plain -t $(NAME) .

local-up: 
	@docker-compose -f "docker/local/docker-compose.yml" up -d --build

local-down: 
	@docker-compose -f "docker/local/docker-compose.yml" down --remove-orphans

wait:
	@sleep 4
		
run:
	@go run cmd/api/*.go

run-watch:
	@air

mig-create: 
	@goose -dir ./migrations create $(MIG_NAME) sql 

mig-status: 
	@goose postgres $(DATABASE.WRITER) status

mig-up: 
	@goose -dir ./migrations postgres $(DATABASE.WRITER) up

mig-down: 
	@goose -dir ./migrations postgres $(DATABASE.WRITER) down

mock: 
	@go generate ./...

test:
	@go test -v -coverpkg=./... -coverprofile=coverage.out ./...
	@go tool cover -func coverage.out

test-cover: test
	@go tool cover -html=coverage.out 

open-swagger:
	open http://localhost:8080/swagger/index.html

open-jaeger:
	open http://localhost:16686/search