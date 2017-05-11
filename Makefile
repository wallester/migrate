
build: fmt
	@go build

fmt:
	@go fmt ./...

lint:
	@gometalinter ./... --vendor

test:
	@go list ./... | grep -v vendor | xargs go test
