# Build configuration
BINARY_NAME=gin-sqlite-demo
MAIN_PATH=./cmd/server
BUILD_DIR=./build
DIST_DIR=./dist

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Build flags
LDFLAGS=-ldflags "-w -s"
BUILD_FLAGS=-trimpath

.PHONY: all build clean test deps tidy run docker

all: clean deps build

build:
	CGO_ENABLED=0 $(GOBUILD) $(BUILD_FLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)

clean:
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)/$(BINARY_NAME)
	rm -rf $(DIST_DIR)

test:
	$(GOTEST) -v ./...

deps:
	$(GOMOD) download

tidy:
	$(GOMOD) tidy

run: build
	$(BUILD_DIR)/$(BINARY_NAME)

docker:
	docker build -t $(BINARY_NAME) -f ../Dockerfile ..

help:
	@echo "Available targets:"
	@echo "  build    - Build the application"
	@echo "  clean    - Clean build artifacts"
	@echo "  test     - Run tests"
	@echo "  deps     - Download dependencies"
	@echo "  tidy     - Tidy go modules"
	@echo "  run      - Build and run the application"
	@echo "  docker   - Build Docker image"
	@echo "  help     - Show this help"
