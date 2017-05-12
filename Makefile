
build: fmt
	@go build

fmt:
	@go fmt ./...

lint:
	@gometalinter ./... --deadline=5m --vendor --enable misspell --enable goimports

test:
	@go list ./... | grep -v vendor | xargs go test

install:
	@go install

cov:
	@go test -coverprofile=coverage.out
    @go tool cover -html=coverage.out
