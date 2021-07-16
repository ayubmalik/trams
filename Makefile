VERSION=$(shell git describe --tags)
LDFLAGS=-s -w
PREVIOUSTAG:=$(shell git describe --tags --abbrev=0)
PREVIOUSTAGDATE:=$(shell git log  -1 --format=%as $(PREVIOUSTAG))
TODAYDATE:=$(shell date +'%Y-%m-%d')

build: test
	go build -ldflags '$(LDFLAGS) -X "main.version=$(VERSION)"' ./cmd/trams/

check-env:

test: clean
	go test ./...	

tag-release: changelog	
	@git add CHANGELOG.md
	@git commit -m  "Release $(NEWTAG) ($(TODAYDATE))"
	@git tag $(NEWTAG)

changelog:
	@echo Previous tag = $(PREVIOUSTAG)
	@echo Previous tag date = $(PREVIOUSTAGDATE)
	@rm -rf change*.tmp
	@echo "# Changelog" > changes.tmp
	@echo "" >> changes.tmp
	@echo "## $(NEWTAG) ($(TODAYDATE))" >> changes.tmp
	@git log $(PREVIOUSTAG)..HEAD --pretty=format:"%h %s" >> changes.tmp
	@echo "" >> changes.tmp
	@sed '/# Changelog/d' CHANGELOG.md > changelog.tmp
	@cat changes.tmp changelog.tmp > newchangelog.tmp
	@mv newchangelog.tmp CHANGELOG.md

clean:
	@go clean -testcache

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'