BINARY=hello
LDFLAGS="-s -w -X main.version=${VERSION}"
PREVIOUSTAG:=$(shell git describe --tags --abbrev=0)
PREVIOUSTAGDATE:=$(shell git log  -1 --format=%as $(PREVIOUSTAG))
TMPCHANGES=/tmp/changes.tmp
TODAYDATE:=$(shell date +'%Y-%m-%d')
VERSION=$(shell git describe --tags --long)

test: clean
	go test ./...

build: test
	go build -ldflags $(LDFLAGS) -o $(BINARY) .

binaries: test
	@mkdir -p dist/linux dist/darwin dist/windows
	GOOS=linux GOARCH=amd64 go build -ldflags $(LDFLAGS) -o dist/linux/$(BINARY) .
	GOOS=darwin GOARCH=amd64 go build -ldflags $(LDFLAGS) -o dist/darwin/$(BINARY) .
	GOOS=windows GOARCH=amd64 go build -ldflags $(LDFLAGS) -o dist/windows/$(BINARY).exe .

changelog: clean
	@echo Previous tag = $(PREVIOUSTAG)
	@echo Previous tag date = $(PREVIOUSTAGDATE)
	@echo "# Changelog" > $(TMPCHANGES)
	@echo "" >> $(TMPCHANGES)
	@echo "## $(NEWTAG) ($(TODAYDATE))" >> $(TMPCHANGES)
	@git log $(PREVIOUSTAG)..HEAD --pretty=format:"%h %s" | grep -v 'build:' | grep -v 'Release v' >> $(TMPCHANGES)
	@echo "" >> $(TMPCHANGES)
	@touch CHANGELOG.md
	@sed '/# Changelog/d' CHANGELOG.md >> $(TMPCHANGES)
	@mv $(TMPCHANGES) CHANGELOG.md

release: changelog
ifndef NEWTAG
	$(error Please set NEWTAG value first, e.g make release NEWTAG=v0.1.)
endif
	@git add CHANGELOG.md
	@git commit -m  "Release $(NEWTAG) ($(TODAYDATE))"
	@git tag $(NEWTAG)
	@git push
	@git push --tags

clean:
	@rm -rf $(TMPCHANGES)
	@rm -rf dist
	@go clean -testcache
