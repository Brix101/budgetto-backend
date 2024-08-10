# Variables
GO_BIN = ~/go/bin
APP_NAME = budgetto
SRC_DIR = ./cmd/budgetto
AIR = air 
GOOSE = goose 
MIGRATIONS_DIR = migrations  

include .env
export

.PHONY: dev
dev:
	${AIR} api

.PHONY: build
build:
	go build -o $(APP_NAME) $(SRC_DIR)

.PHONY: migrate-up
migrate-up:
	cd ${MIGRATIONS_DIR} && $(GOOSE) postgres $(DATABASE_URL) up

.PHONY: migrate-down
migrate-down:
	cd ${MIGRATIONS_DIR} && $(GOOSE) postgres $(DATABASE_URL) down

.PHONY: migrate-reset
migrate-reset:
	cd ${MIGRATIONS_DIR} && $(GOOSE) postgres $(DATABASE_URL) reset

.PHONY: migrate-status
migrate-status:
	cd ${MIGRATIONS_DIR} && $(GOOSE) postgres $(DATABASE_URL) status
