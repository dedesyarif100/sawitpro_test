

.PHONY: clean all init generate generate_mocks migrate rollback

all: build/main

build/main: cmd/main.go generated
	@echo "Building..."
	go build -o $@ $<

clean:
	rm -rf generated

init: clean generate
	go mod tidy
	go mod vendor

test:
	go clean -testcache
	go test -short -coverprofile coverage.out -short -v ./...

coverage:
	go clean -testcache
	go test -coverprofile=coverage.out -covermode=atomic ./...

coverage-html:
	go clean -testcache
	go test -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -html=coverage.out

coverage-func:
	go clean -testcache
	go test -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -func=coverage.out

test_api:
	go clean -testcache
	go test ./tests/...

migrate:
	@echo "Applying database.sql to the running Postgres container..."
	docker compose exec -T db psql -v ON_ERROR_STOP=1 -U postgres -d database < database.sql

rollback:
	@echo "Rolling back dummy data from the running Postgres container..."
	docker compose exec -T db psql -v ON_ERROR_STOP=1 -U postgres -d database < database.down.sql

generate: generated generate_mocks

generated: api.yml
	@echo "Generating files..."
	mkdir generated || true
	oapi-codegen --package generated -generate types,server,spec $< > generated/api.gen.go

INTERFACES_GO_FILES := $(shell find repository -name "interfaces.go")
INTERFACES_GEN_GO_FILES := $(INTERFACES_GO_FILES:%.go=%.mock.gen.go)

generate_mocks: $(INTERFACES_GEN_GO_FILES)
$(INTERFACES_GEN_GO_FILES): %.mock.gen.go: %.go
	@echo "Generating mocks $@ for $<"
	mockgen -source=$< -destination=$@ -package=$(shell basename $(dir $<))