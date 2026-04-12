projectname?=speedtest-cli

default: help

.PHONY: help
help: ## list makefile targets
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: build
build: ## build golang binary
	mkdir -p dist
	@go build -ldflags "-X main.version=$(shell git describe --abbrev=0 --tags)" -o dist/$(projectname) cmd/cli/main.go

.PHONY: run
run: ## run the app
	@go run -ldflags "-X main.version=$(shell git describe --abbrev=0 --tags)"  cmd/cli/main.go

PHONY: test
test: clean ## display test coverage
	go test --cover -v -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out | sort -rnk3
	
PHONY: clean
clean: ## clean up environment
	@rm -rf coverage.out dist/

PHONY: cover
race: ## display test coverage with race
	go test -v -race $(shell go list ./... | grep -v /vendor/) -v -coverprofile=coverage.out
	go tool cover -func=coverage.out

PHONY: snapshot
snapshot: ## goreleaser snapshot
	goreleaser release --snapshot --clean
