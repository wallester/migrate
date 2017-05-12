
build: fmt
	@go build

fmt:
	@go fmt ./...

lint:
	@gometalinter ./... --vendor

test:
	@go list ./... | grep -v vendor | xargs go test

install:
	@go install

cov:
	@go test -coverprofile=coverage.out
    @go tool cover -html=coverage.out
