export GO111MODULE=on
export NOCACHE := $(if $(NOCACHE),"--no-cache")
export SERVICE="WAS"
export SERVICE_FULLNAME="service-was"

GOOSE_MASTER_DSN="postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"

# OpenAPI code generation
OAPI_CFG := oapi.cfg.yaml
OAPI_SPEC := api/service/openapi.yaml

######################### COMMON #########################

.PHONY: run
run:
	@go run ./cmd/main.go

.PHONY: build
build:
	@go build -o bin/service cmd/main.go

.PHONY: gen
gen: clean-gen gen-sql gen-oapi gen-go

.PHONY: clean-gen
clean-gen:
	@rm -rf ./internal/generated

.PHONY: gen-go
gen-go:
	@go generate ./...

.PHONY: gen-oapi
gen-oapi:
	@mkdir -p internal/generated/service
	@oapi-codegen --config api/service/oapi-service.cfg.yaml $(OAPI_SPEC)

.PHONY: gen-sql
gen-sql:
	@go run -mod=mod github.com/sqlc-dev/sqlc/cmd/sqlc@v1.29.0 generate

.PHONY: gen-test
gen-test:
	./bin/testgen -path "$(CURDIR)/$(PATH)"

.PHONY: tools
tools:
	@go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v2.5.0
	@go install github.com/pressly/goose/v3/cmd/goose@v3.25.0

.PHONY: migrations-apply
migrations-apply:
	@goose --dir ./db/migrations/ postgres "${GOOSE_MASTER_DSN}" up

.PHONY: migrations-down
migrations-down:
	@goose --dir ./db/migrations/ postgres "${GOOSE_MASTER_DSN}" down

.PHONY: migrations-create
migrations-create:
	@goose --dir ./db/migrations/ create $(NAME) sql

.PHONY: lint
lint:
	@golangci-lint run
