
export GOBIN := $(PWD)/bin
export PATH := $(GOBIN):$(PATH)

TOOLS=$(shell cat tools/tools.go | egrep '^\s_ '  | awk '{ print $$2 }')

.PHONY: bootstrap-tools
bootstrap-tools:
	@echo "Installing: " $(TOOLS)
	@cd tools && go install $(TOOLS)

.PHONY: run
run:
	go run main.go $(ARGS)

.PHONY: lint
lint:
	$(GOBIN)/golangci-lint run -v ./...

.PHONY: lint-fix
lint-fix:
	$(GOBIN)/golangci-lint run --fix -v ./...

.PHONY: test
test:
	go test -v -race ./...

.PHONY: cover
cover:
	go test -v -race -coverpkg=./... -coverprofile=coverage.out ./...

.PHONY: tidy
tidy:
	go mod tidy
	cd tools && go mod tidy
