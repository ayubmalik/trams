LDFLAGS="-s -w"

build: test
	go build -ldflags=$(LDFLAGS) ./cmd/trams/
 
check-env:

test: clean
	go test ./
 
clean:
	@go clean -testcache
