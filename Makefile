-include .env.local
current_dir := $(patsubst %/,%,$(dir $(abspath $(firstword $(MAKEFILE_LIST)))))

export TOOLS=$(current_dir)/tools
export TOOLS_BIN=$(TOOLS)/bin
export PATH := $(TOOLS_BIN):$(PATH)

.PHONY: 
run:
	docker-compose up -d
	go run cmd/main.go

.PHONY:
fix-lint: install-tools
	golangci-lint run --fix

.PHONY:
lint: install-tools
	golangci-lint run

.PHONY:
install-tools: export GOBIN=$(TOOLS_BIN)
install-tools:
	@if [ -d "$(TOOLS_BIN)" ] && [ -n "$$(ls -A $(TOOLS_BIN))" ]; then \
		echo "TOOLS_BIN directory exists and not empty"; \
	else \
		echo "TOOLS_BIN directory does not exist or is empty.\n Installing tools..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.61.0; \
	fi

.PHONY:
setup-pre-commit:
	@chmod +x scripts/setup-pre-commit.sh
	@./scripts/setup-pre-commit.sh
