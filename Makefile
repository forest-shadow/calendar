-include .env.local
current_dir := $(patsubst %/,%,$(dir $(abspath $(firstword $(MAKEFILE_LIST)))))

export TOOLS=$(current_dir)/tools
export TOOLS_BIN=$(TOOLS)/bin
export PATH := $(TOOLS_BIN):$(PATH)

.PHONY:
.SILENT: 
run:
	docker-compose up postgres -d
	echo "Waiting for PostgreSQL to be ready..."
	until docker-compose exec -T postgres pg_isready -U postgres; do \
		sleep 1; \
	done
	echo "PostgreSQL is ready"
	go run cmd/main.go

.PHONY:
run-docker:
	docker-compose up -d

.PHONY:
fix-lint: install-tools
	golangci-lint run --fix

.PHONY:
lint: install-tools
	golangci-lint run

.PHONY:
.SILENT:
install-tools: export GOBIN=$(TOOLS_BIN)
install-tools:
	if [ -d "$(TOOLS_BIN)" ] && [ -n "$$(ls -A $(TOOLS_BIN))" ]; then \
		echo "TOOLS_BIN directory exists and not empty"; \
	else \
		echo "TOOLS_BIN directory does not exist or is empty.\n Installing tools..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.61.0; \
	fi

.PHONY:
.SILENT:
setup-pre-commit:
	chmod +x scripts/setup-pre-commit.sh
	./scripts/setup-pre-commit.sh
