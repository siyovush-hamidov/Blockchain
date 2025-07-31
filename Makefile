# Default target
.PHONY: all
all: build

# Build the node and client executables
.PHONY: build
build:
	go get ./...
	go build -o node node.go
	go build -o client client.go

# Remove compiled binaries
.PHONY: clean
clean:
	rm -f node client

# Display available commands
.PHONY: help
help:
	@echo "Available commands:"
	@echo "  make build  - Compile node and client executables"
	@echo "  make clean  - Remove compiled binaries"
	@echo "  make help   - Show this help message"