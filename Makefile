VERSION=$(shell git describe --tags)
LDFLAGS=-s -w
PREVIOUSTAG:=$(shell git describe --tags --abbrev=0)
PREVIOUSTAGDATE:=$(shell git log  -1 --format=%as $(PREVIOUSTAG))
TODAYDATE:=$(shell date +'%Y-%m-%d')
TMPCHANGES=$(TMPDIR)/changes.tmp

build: test
	go build -ldflags '$(LDFLAGS) -X "main.version=$(VERSION)"' ./cmd/trams/

check-env:

test: clean
	go test ./...

tag-release: changelog
	@git add CHANGELOG.md
	@git commit -m  "Release $(NEWTAG) ($(TODAYDATE))"
	@git tag $(NEWTAG)

changelog: clean
	@echo Previous tag = $(PREVIOUSTAG)
	@echo Previous tag date = $(PREVIOUSTAGDATE)
	@echo "# Changelog" > $(TMPCHANGES)
	@echo "" >> $(TMPCHANGES)
	@echo "## $(NEWTAG) ($(TODAYDATE))" >> $(TMPCHANGES)
	@git log $(PREVIOUSTAG)..HEAD --pretty=format:"%h %s" >> $(TMPCHANGES)
	@echo "" >> $(TMPCHANGES)
	@sed '/# Changelog/d' CHANGELOG.md >> $(TMPCHANGES)
	@cat $(TMPCHANGES)


clean:
	@rm -rf change*.tmp
	@go clean -testcache

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'