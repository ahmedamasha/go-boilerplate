BINARY_NAME=refda-api
BUILD_DIR=build
GOCMD=go

.PHONY: run build test test-api clean deps docker-up docker-down vendor web-dev web-build help

run:
	$(GOCMD) run ./cmd/api

build:
	@mkdir -p $(BUILD_DIR)
	$(GOCMD) build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/api

test:
	$(GOCMD) test ./...

test-api:
	./scripts/test_api.sh

clean:
	$(GOCMD) clean
	rm -rf $(BUILD_DIR) uploads

deps:
	$(GOCMD) mod download
	$(GOCMD) mod tidy

docker-up: vendor
	docker compose up --build -d

vendor:
	$(GOCMD) mod vendor

docker-down:
	docker compose down

WEB_NPM = ./scripts/web_npm.sh

web-dev:
	$(WEB_NPM) run dev

web-build:
	$(WEB_NPM) run build

help:
	@echo "run, build, test, clean, deps, vendor, docker-up, docker-down, web-dev, web-build"
