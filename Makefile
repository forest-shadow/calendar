YAML_FILE=env.local.yml

ifeq ($(shell which yq),)
  $(error "yq is not installed. Please install yq to use this Makefile.")
endif

HTTP_PORT=$(shell yq eval '.http.port' $(YAML_FILE))
DB_URI=$(shell yq eval '.db.uri' $(YAML_FILE))

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
test: install-tools
	hurl --test --verbose --variable host=localhost:$(HTTP_PORT) scripts/requests/$(name).hurl

.PHONY:
migration-create: install-tools
	$(TOOLS_BIN)/goose create -dir ./migrations "$(name)" sql

.PHONY:
migrate: install-tools
	$(TOOLS_BIN)/goose -dir ./migrations postgres "$(DB_URI)" up -v

.PHONY:
.SILENT:
install-tools: export GOBIN=$(TOOLS_BIN)
install-tools:
	if [ -d "$(TOOLS_BIN)" ] && [ -n "$$(ls -A $(TOOLS_BIN))" ]; then \
		echo "TOOLS_BIN directory exists and not empty"; \
	else \
		echo "TOOLS_BIN directory does not exist or is empty.\n Installing tools..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.61.0; \
		go install github.com/pressly/goose/v3/cmd/goose@v3.22.1; \
	fi

.PHONY:
.SILENT:
setup-pre-commit:
	chmod +x scripts/setup-pre-commit.sh
	./scripts/setup-pre-commit.sh
