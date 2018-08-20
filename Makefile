
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

tools:
	@echo "govendor" && go get -u github.com/kardianos/govendor
	@echo "golangci-lint" && go get -u github.com/golangci/golangci-lint/cmd/golangci-lint
