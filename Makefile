export GOFLAGS=-mod=vendor

build: fmt
	@go build

fmt:
	@go fmt ./...

lint: install
	@golangci-lint run

test:
	@go list ./... | xargs go test

install:
	@go install ./...

cov:
	@go test -test.covermode=count -test.coverprofile coverage.cov
	@go tool cover -html=coverage.cov -o coverage.html

validations: test lint