GO=GO111MODULE=on go
GO_OFF=GO111MODULE=off go
GOLINT=golangci-lint

.PHONY: all
all: lint test bins


.PHONY: bins
bins:
	$(GO) build -o ./bin/gin_project_starter ./src

.PHONY: lint
lint:
	$(GOLINT) run --fix --skip-dirs vendor
	$(GO) vet ./...
	$(GO) fmt ./...

.PHONY: test
test:
	$(GO) test -cover -race -v ./...
