# console foreground colors
FG_RED := `tput setaf 1`
FG_GREEN := `tput setaf 2`
FG_RESET := `tput sgr0`

# PostgreSQL parameters
PG_USER := postgres
PG_PASSWORD := mysecret
PG_HOST := localhost
PG_PORT := 30010
DBNAME := protection
MIGRATIONS_DIR := ./db/migrations

.PHONY: check-go-version
check-go-version:
	@echo "$(FG_GREEN)Checking Go version$(FG_RESET)"
	$(eval GO_MINOR_VERSION := $(shell go version | cut -d " " -f 3 | cut -d "." -f 2))
	@if [ $(GO_MINOR_VERSION) -lt 18 ]; then \
		$(bold) \
		echo "$(FG_RED)Go version 1.18 or higher is required$(FG_RESET)"; \
		exit 1; \
	fi

.PHONY: install-tools
install-tools: check-go-version
	@echo "$(FG_GREEN)Installing migrate$(FG_RESET)"
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@echo "$(FG_GREEN)Installing sqlc$(FG_RESET)"
	go install github.com/kyleconroy/sqlc/cmd/sqlc@latest

.PHONY: migrate-up
migrate-up:
	@echo "$(FG_GREEN)Running DB migrations$(FG_RESET)"
	migrate -path $(MIGRATIONS_DIR) -database "postgresql://$(PG_USER):$(PG_PASSWORD)@$(PG_HOST):$(PG_PORT)/$(DBNAME)?sslmode=disable" -verbose up
