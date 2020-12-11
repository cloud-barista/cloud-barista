SHELL := /bin/bash

GOCMD=go
GOTEST=$(GOCMD) test
GOCLEAN=$(GOCMD) clean
GOBUILD=$(GOCMD) build
GOINSTALL=$(GOCMD) install

BINARY_NAME=cb-restapigw
BINARY_LINUX=$(BINARY_NAME)
TAGS=$(BINARY_NAME)
CONF=./conf/$(BINARY_NAME).yaml

export CBLOG_ROOT=$(CURDIR)
export CBSTORE_ROOT=$(CURDIR)

.PHONY: all clean build
all: clean build

.PHONY: clean
clean:
	@echo "==> Cleaning project"
	rm -f $(BINARY_NAME)

.PHONY: build-run build-run-linux run check
build-run: build-dev run
build-run-linux: build-dev-linux run
build-check: build-dev check

check:
	@echo "==> Check the configuration file"
	./$(BINARY_NAME) --config $(CONF) check -d true -p 8080

run:
	@echo "==> Build project"
	./$(BINARY_NAME) --config $(CONF) run

.PHONY: build build-dev build-dev-linux
build: static
	@echo "==> Build project"
	$(GOBUILD) -tags "$(TAGS)" -o $(BINARY_NAME) -v $(BUILD_PATH)

build-dev:
	@echo "==> Build project (Dev)"
	$(GOBUILD) -tags "$(TAGS)" -o $(BINARY_NAME) -v $(BUILD_PATH)

build-dev-linux:
	@echo "==> Build (Linux-64) project (Dev)"
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -tags "$(TAGS)" -o $(BINARY_NAME) -v $(BUILD_PATH)

.PHONY: build-linux build-check-linux
build-linux: static
	@echo "==> Build (Linux-64) project"
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -tags "$(TAGS)" -o $(BINARY_LINUX) -v $(BUILD_PATH)

build-check-linux: build-dev-linux check

.PHONY: install
install:
	@echo "==> Install project"
	$(GOINSTALL) -tags "$(TAGS)"

.PHONY: static clean-client client 
static: client
	esc -o ${PWD}/pkg/static/static.go -ignore "DS_Store|.nojekyll|README.md" -pkg static web/dist

clean-client:
	rm -rf ${PWD}/web/.nuxt ${PWD}/web/dist

client: clean-client
	cd ${PWD}/web; yarn generate