# Define variables
GO      := go
GOBUILD := $(GO) build
GOTEST  := $(GO) test
GOCLEAN := $(GO) clean
BINARY  := postgres-crud

.PHONY: all build test clean run
	
all: build

build: ## Compile the application binary
	$(GOBUILD) -o $(BINARY) .

test: ## Run tests
	$(GOTEST) -v ./...

clean: ## Remove the compiled binary
	$(GOCLEAN)
	rm -f $(BINARY)

run: build ## Build and run the application
	./$(BINARY)

help: ## Show this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'