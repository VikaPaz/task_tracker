POSTGRES_PASSWORD ?= password
POSTGRES_USER ?= user
POSTGRES_DB ?= tasks

run:
	go run cmd/main.go

build: 
	docker-compose -f build/docker-compose.yaml up --build

run_postgres:
	POSTGRES_PASSWORD=${POSTGRES_PASSWORD} \
	POSTGRES_USER=${POSTGRES_USER} \
	POSTGRES_DB=${POSTGRES_DB} \
	docker-compose -f build/docker-compose.yaml up postgres
