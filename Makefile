GO=GO111MODULE=on go
GOLINT=golangci-lint

.PHONY: all
all: lint test build


.PHONY: build
build:
	$(GO) build -o ./bin/gin_project_starter ./src

.PHONY: run
run:
	$(GO) run ./src

.PHONY: lint
lint:
	$(GOLINT) run --fix --skip-dirs vendor
	$(GO) vet ./...
	$(GO) fmt ./...

.PHONY: test
test:
	$(GO) test -coverprofile cover.out -race -v ./...
