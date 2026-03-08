.DEFAULT_GOAL := all

BUILD = $(shell git rev-parse HEAD 2> /dev/null || echo "undefined")
BUILDDATE = $(shell LANG=en_us_88591 date)
CGO_ENABLED := 0
DEBUG := 0
UNAME_m := $(shell uname -m)
UNAME_s := $(shell uname -s)
VERSION = $(shell git describe --abbrev=0 --tags 2> /dev/null || echo "0.1.0")

# Check if it's a development build by verifying the DEBUG flag;
# should that be any value different from '0' strip and trim the binary
ifeq ($(DEBUG),0)
	GO_LDFLAGS_client := -ldflags '-s -w -extldflags "-static"'
	# If we are building on top of a macOS/arm64, avoid setting the -static flag;
	# usually additional steps and softwares would be necessary for the setup to work (ld: crt0.o).
    ifeq ($(UNAME_s),Darwin)
	    ifeq ($(UNAME_m),arm64)
		    GO_LDFLAGS_server := -ldflags '-s -w'
		else
		    GO_LDFLAGS_server := -ldflags '-s -w -extldflags "-static"'
	    endif # if UNAME_m=arm64
    else
		GO_LDFLAGS_server := -ldflags '-s -w -extldflags "-static"'
    endif # if UNAME_s=Darwin
else
    GO_LDFLAGS_client :=
    GO_LDFLAGS_server :=
endif # if DEBUG=0

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: all
all: ## Build both the client and the server in their respective directories
	go build -o ./client/client -trimpath $(GO_LDFLAGS_client) ./client
	go build -o ./server/server -trimpath $(GO_LDFLAGS_server) ./server

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
	-rm -r goploader-server
	-rm -r gpldr
	-rm -r releases/
