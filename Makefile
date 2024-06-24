-include .env
export

CURRENT_DIR=$(shell pwd)

# run service
.PHONY: run
run:
	go run cmd/app/main.go

swagger-gen:
	~/go/bin/swag init -g ./api/router.go -o api/docs


# go generate
.PHONY: go-gen
go-gen:
	go generate ./...

# migrate
.PHONY: migrate
migrate:
	migrate -source file://migrations -database postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DATABASE}?sslmode=disable up

DB_URL := "postgres://postgres:123@localhost:5432/generator_resume?sslmode=disable"

migrate-up:
	migrate -path migrations -database $(DB_URL) -verbose up

migrate-down:
	migrate -path migrations -database $(DB_URL) -verbose down

migrate-force:
	migrate -path migrations -database $(DB_URL) -verbose force 1

migrate-file:
	migrate create -ext sql -dir migrations/ -seq create_soft_skills_table

