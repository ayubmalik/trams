VERSION=$(shell git describe --tags)
LDFLAGS=-s -w

build: test
	go build -ldflags '$(LDFLAGS) -X "main.version=$(VERSION)"' ./cmd/trams/

check-env: 		## check any required environment variables

test: clean 	## clean go testcache
	go test ./...	

tag-release: 	## tag git repo with specified RELEASE and update CHANGELOG.md
	echo release is $${RELEASE}
	@CHANGES=$(git log $(git describe --tags --abbrev=0)..HEAD --pretty=format:"%h %s")
	echo $${CHANGES}
clean:
	@go clean -testcache

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'