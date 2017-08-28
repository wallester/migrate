
build: fmt
	@go build

fmt:
	@go fmt ./...

lint:
	@gometalinter.v1 ./... --config=.gometalinter

test:
	@go list ./... | grep -v vendor | xargs go test

install:
	@go install

cov:
	@go test -test.covermode=count -test.coverprofile coverage.cov
	@go tool cover -html=coverage.cov -o coverage.html

tools:
	@echo "govendor" && go get -u github.com/kardianos/govendor
	@echo "gometalinter.v1" && go get -u gopkg.in/alecthomas/gometalinter.v1
