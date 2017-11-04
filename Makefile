
build: fmt
	@go build

fmt:
	@go fmt ./...

lint: install
	@gometalinter ./... --config=.gometalinter

test:
	@go list ./... | xargs go test

install:
	@go install ./...

cov:
	@go test -test.covermode=count -test.coverprofile coverage.cov
	@go tool cover -html=coverage.cov -o coverage.html

tools:
	@echo "govendor" && go get -u github.com/kardianos/govendor
	@echo "gometalinter" && go get -u github.com/alecthomas/gometalinter
