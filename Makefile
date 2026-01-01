.PHONY: help format format-check lint install-tools check-gofumpt check-golangci-lint

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

format: check-gofumpt ## Format code with gofumpt
	gofumpt -w .

format-check: check-gofumpt ## Check if code is properly formatted
	@echo "Checking formatting..."
	@test -z "$$(gofumpt -l .)" || { \
		echo "Code is not properly formatted. Run 'make format' to fix."; \
		exit 1; \
	}
	@echo "Code is properly formatted."

lint: check-golangci-lint ## Run golangci-lint
	golangci-lint run

install-tools: ## Install gofumpt and golangci-lint
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
