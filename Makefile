
build: fmt
	@go build

fmt:
	@go fmt ./...

lint:
	@gometalinter ./... --deadline=5m --vendor --enable misspell --enable goimports --disable aligncheck --concurrency 4

test:
	@go list ./... | grep -v vendor | xargs go test

install:
	@go install

cov:
	@go test -test.covermode=count -test.coverprofile coverage.cov
	@go tool cover -html=coverage.cov -o coverage.html
