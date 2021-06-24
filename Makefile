LDFLAGS="-s -w"

build: test
	go build -ldflags=$(LDFLAGS) ./cmd/trams/
 
check-env:

test: clean
	go test ./

build:
	go build -ldflags=$(LDFLAGS) ./cmd/trams
 
clean:
	@go clean -testcache
