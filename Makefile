.PHONY: all init-db seed-db run

# Variables
include .env
GO := go
GOFLAGS :=

# Targets
all: init-db seed-db run

init-db:
	DB_DSN=$(DB_DSN) go run cmd/tools/initdb/main.go

seed-db:
	DB_DSN=$(DB_DSN) go run cmd/tools/seeddb/main.go

run:
	DB_DSN=$(DB_DSN) go run cmd/server/main.go

# Test target to run all tests
.PHONY: test
test:
	go test -v ./...

# Clean target to remove any build artifacts (optional)
.PHONY: clean
clean:
	$(GO) clean -testcache
