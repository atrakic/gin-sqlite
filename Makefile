# Go parameters
GOCMD = go
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean
GOTEST = $(GOCMD) test
GOGET = $(GOCMD) get

# Main package name
MAIN_PACKAGE = github.com/atrakic/gin-sqlite

# Output binary name
BINARY_NAME = server

# Build the project
build:
	$(GOBUILD) -o $(BINARY_NAME) $(MAIN_PACKAGE)

# Clean the project
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

# Run tests
test:
	$(GOTEST) -v ./...

# Install project dependencies
deps:
	$(GOGET) -v ./...

# Run the project
run:
	$(GOBUILD) -o $(BINARY_NAME) $(MAIN_PACKAGE)
	./$(BINARY_NAME)

# Default target
default: build
