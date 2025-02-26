.DEFAULT_GOAL := all
CGO_ENABLED=0
VERSION=$(shell git describe --abbrev=0 --tags 2> /dev/null || echo "0.1.0")
BUILD=$(shell git rev-parse HEAD 2> /dev/null || echo "undefined")
BUILDDATE=$(shell LANG=en_us_88591 date)

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: all
all: ## Build both the client and the server in their respective directories
	go build -o ./client/client ./client
	go build -o ./server/server ./server

.PHONY: docker
docker: ## Build the docker image
	docker build -t gpldr:latest -t gpldr:$(BUILD) -f Dockerfile .

.PHONY: release
release: ## Create a new release
	goreleaser release --clean

.PHONY: snapshot
snapshot: ## Create a new snapshot release
	goreleaser release --snapshot --clean

clean:
	-rm -r releases/
	-rm -r goploader-server
