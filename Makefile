
build: fmt
	@go build

fmt:
	@go fmt ./...

lint:
	@gometalinter ./...

test:
	@go test ./...