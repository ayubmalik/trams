VERSION=$(shell git describe --tags)
LDFLAGS=-s -w

build: test
	go build -ldflags '$(LDFLAGS) -X "main.version=$(VERSION)"' ./cmd/trams/

check-env:

test: clean
	go test ./

clean:
	@go clean -testcache
