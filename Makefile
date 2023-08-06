help:
	@echo "Please use \`make <target>' where <target> is one of:"
	@echo "  build         Compiles and builds the application."
	@echo "  dev           Runs the application and a postgres database via Docker compose."
	@echo "  test          Run unit tests."

all: build test

build:
	go build -v ./...

dev:
	docker compose up

test:
	go test -v ./... -coverprofile=coverage.txt -covermode=atomic


.PHONY: help build dev test