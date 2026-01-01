.PHONY: help format format-check lint test all install-tools check-gofumpt check-golangci-lint

help:
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

format: check-gofumpt
	gofumpt -w .

format-check: check-gofumpt
	@echo "Checking formatting..."
	@test -z "$$(gofumpt -l .)" || { \
		echo "Code is not properly formatted. Run 'make format' to fix."; \
		exit 1; \
	}
	@echo "Code is properly formatted."

lint: check-golangci-lint
	golangci-lint run

test:
	go test ./...

all: format-check lint test

install-tools:
	go install mvdan.cc/gofumpt@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

check-gofumpt:
	@command -v gofumpt >/dev/null 2>&1 || { \
		echo "gofumpt não está instalado. Execute: make install-tools"; \
		exit 1; \
	}

check-golangci-lint:
	@command -v golangci-lint >/dev/null 2>&1 || { \
		echo "golangci-lint não está instalado. Execute: make install-tools"; \
		exit 1; \
	}
