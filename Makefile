VERSION=$(shell git describe --tags)
LDFLAGS=-s -w
PREVIOUSTAG:=$(shell git describe --tags --abbrev=0)

build: test
	go build -ldflags '$(LDFLAGS) -X "main.version=$(VERSION)"' ./cmd/trams/

check-env:

test: clean
	go test ./...	

tag-release: changelog	
	@echo Previous release tag = $(PREVIOUSTAG)
	@echo New release tag = $(RELEASETAG)

changelog:
	@echo Previous tag = $(PREVIOUSTAG)
	git log $(PREVIOUSTAG)..HEAD --pretty=format:"%h %s" > changes.txt

clean:
	@go clean -testcache

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'