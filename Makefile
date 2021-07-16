VERSION=$(shell git describe --tags)
LDFLAGS=-s -w

build: test
	go build -ldflags '$(LDFLAGS) -X "main.version=$(VERSION)"' ./cmd/trams/

check-env:

test: clean
	go test ./...	
tag-release:
	echo release is $${RELEASE}
	@CHANGES=$(git log $(git describe --tags --abbrev=0)..HEAD --pretty=format:"%h %s")
	echo $${CHANGES}
clean:
	@go clean -testcache
