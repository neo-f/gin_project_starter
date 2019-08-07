.PHONY: all
all: lint test build

.PHONY: build
build:
	go build -o ./bin/gin_project_starter ./src

.PHONY: lint
lint:
	golangci-lint run -v --fix --skip-dirs vendor

.PHONY: run
run:
	go run ./src

.PHONY: test
test:
	go test -coverprofile=coverage.txt -race -v ./...

.PHONY: setup
setup:
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s -- -b $(GOPATH)/bin
	golangci-lint --version
