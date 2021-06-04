
export GOBIN := $(PWD)/bin
export PATH := $(GOBIN):$(PATH)

TOOLS=$(shell cat tools/tools.go | egrep '^\s_ '  | awk '{ print $$2 }')

.PHONY: bootstrap-tools
bootstrap-tools:
	@echo "Installing: " $(TOOLS)
	@go install $(TOOLS)

.PHONY: run-example
run-example:
	go run cmd/yaml-diff/main.go -file1 example/a.yaml -file2 example/b.yaml
	@echo --------------------
	go run cmd/yaml-diff/main.go -file1 example/b.yaml -file2 example/a.yaml

.PHONY: lint
lint:
	golangci-lint run -v ./...
	go-consistent -v ./...

.PHONY: lint-fix
lint-fix:
	golangci-lint run --fix -v ./...

.PHONY: test
test:
	go test -v -race ./...

.PHONY: cover
cover:
	go test -v -race -coverpkg=./... -coverprofile=coverage.txt ./...

.PHONY: tidy
tidy:
	go mod tidy
