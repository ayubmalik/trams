VERSION=$(shell git describe --tags)
LDFLAGS=-s -w
PREVIOUSTAG:=$(shell git describe --tags --abbrev=0)
PREVIOUSTAGDATE:=$(shell git log  -1 --format=%as $(PREVIOUSTAG))
TODAYDATE:=$(shell date +'%Y-%m-%d')
TMPCHANGES=$(TMPDIR)/changes.tmp
LDFLAGS="-s -w"

build: test
	go build -ldflags '$(LDFLAGS) -X "main.version=$(VERSION)"' ./cmd/trams/

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
	@touch CHANGELOG.md
	@sed '/# Changelog/d' CHANGELOG.md >> $(TMPCHANGES)
	@mv $(TMPCHANGES) CHANGELOG.md

binaries: clean
	@mkdir -p dist/linux dist/darwin dist/windows
	GOOS=darwin GOARCH=amd64 go build -ldflags=$(LDFLAGS) -o dist/darwin/trams ./cmd/trams
	
clean:
	@rm -rf $(TMPCHANGES)
	@rm -rf dist
	@go clean -testcache
