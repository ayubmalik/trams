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
	@echo New release tag = $(NEWTAG)

changelog:
	@echo Previous tag = $(PREVIOUSTAG)
	@rm -rf change*.tmp
	git log $(PREVIOUSTAG)..HEAD --pretty=format:"%h %s" > changes.tmp
	echo >> changes.tmp
	sed '/# Changelog/d' CHANGELOG.md > changelog.tmp
	cat changes.tmp changelog.tmp > newchangelog.tmp
clean:
	@go clean -testcache

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'