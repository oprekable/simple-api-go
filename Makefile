.PHONY: docker-up
docker-up:
	@docker compose up -d

.PHONY: docker-down
docker-down:
	@docker compose -p simple-api-go down -v

.PHONY: docker-clean-data
docker-clean-data:
	@rm -rf .docker/db/data

.PHONY: go-mod-tidy
go-mod-tidy:
	@CGO_ENABLED=1 go mod tidy

.PHONY: install-fieldalignment
install-fieldalignment:
	@go install golang.org/x/tools/go/analysis/passes/fieldalignment/cmd/fieldalignment@latest

.PHONY: fix-struct
fix-struct:
	@CGO_ENABLED=1 fieldalignment -fix ./...

.PHONY: install-golangci-lint
install-golangci-lint:
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

.PHONY: go-lint
go-lint:
	@CGO_ENABLED=1 golangci-lint run ./... --fix

.PHONY: go-lint-fix-struct
go-lint-fix-struct:
	@CGO_ENABLED=1 go mod tidy
	@CGO_ENABLED=1 fieldalignment -fix ./...
	@CGO_ENABLED=1 golangci-lint run ./... --fix

.PHONY: install-staticcheck
install-staticcheck:
	@go install honnef.co/go/tools/cmd/staticcheck@latest

.PHONY: go-staticcheck
go-staticcheck:
	@staticcheck ./...

.PHONY: migrate-dev-install
migrate-dev-install:
	@go install github.com/rubenv/sql-migrate/...@latest

.PHONY: migrate-dev-new
migrate-dev-new:
	@read -p "Enter Migrate Name:" migname; \
	sql-migrate new -env local "$$migname";

.PHONY: migrate-dev-up
migrate-dev-up:
	@sql-migrate up -env local

.PHONY: migrate-dev-down
migrate-dev-down:
	@sql-migrate down -env local -limit=0

.PHONY: run-app
run-app:
	@go run . --env=local

.PHONY: test-api-no-filter
test-api-no-filter:
	@curl -s --request GET --url http://localhost:3000/api | jq

.PHONY: test-api-filter-brand
test-api-filter-brand:
	@curl -s --request GET --url 'http://localhost:3000/api?brand=Honda' | jq

.PHONY: test-api-filter-type
test-api-filter-type:
	@curl -s --request GET --url 'http://localhost:3000/api?type=Beat' | jq

.PHONY: test-api-filter-transmission
test-api-filter-transmission:
	@curl -s --request GET --url 'http://localhost:3000/api?transmission=Manual' | jq

.PHONY: test-api-filter-all-fields
test-api-filter-all-fields:
	@curl -s --request GET --url 'http://localhost:3000/api?brand=Honda&type=Beat&transmission=Automatic' | jq
