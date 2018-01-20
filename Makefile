PACKAGES = $(shell go list ./... | grep -v /vendor/)
EXTRA_DEPENDENCIES = \
    github.com/golang/lint/golint \
    github.com/tools/godep \
    github.com/wadey/gocovmerge \
    github.com/alecthomas/gometalinter


# Same as $(PACKAGES) except we get directory paths. We exclude the first line
# because it contains the top level directory which contains /vendor/
PACKAGE_DIRS=$(shell go list -f '{{ .Dir }}' ./... | egrep -v /vendor/ | tail -n +2)

SOURCES = $(shell for f in $(PACKAGES); do ls $$GOPATH/src/$$f/*.go; done)

check: deps vet test coverage lint

deps:
	go get $(EXTRA_DEPENDENCIES)
	gometalinter --install > /dev/null

fmt:
	go fmt ./...
	goimports -w $(SOURCES)

vet:
	go vet $(PACKAGES)

build:
	go build .

test:
	go test -race ./... -v -check.v

# coverage runs the tests to collect coverage but does not attempt to look
# for race conditions.
coverage: $(patsubst %,%.coverage,$(PACKAGES))
	@rm -f .gocoverage/cover.txt
	gocovmerge .gocoverage/*.out > coverage.txt
	go tool cover -html=coverage.txt -o .gocoverage/index.html

%.coverage:
	@[ -d .gocoverage ] || mkdir .gocoverage
	go test -v -covermode=count -coverprofile=.gocoverage/$(subst /,-,$*).out $* -check.v

lint:
	gometalinter --vendor --disable-all --enable=deadcode --enable=errcheck --enable=goimports \
	--enable=gocyclo --enable=golint --enable=gosimple --enable=misspell \
	--enable=unconvert --enable=unused --enable=varcheck --enable=interfacer \
	./...
