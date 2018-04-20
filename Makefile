NAME=zalua
BINARY=./bin/${NAME}
SOURCEDIR=./src
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

VERSION := $(shell git describe --abbrev=0 --tags)
SHA := $(shell git rev-parse --short HEAD)

GOPATH ?= /usr/local/go
GOPATH := ${CURDIR}:${GOPATH}
export GOPATH

$(BINARY): $(SOURCES)
	go build -o ${BINARY} -ldflags "-X main.BuildVersion=$(VERSION)-$(SHA)" $(SOURCEDIR)/$(NAME)/cmd/main.go

run: clean $(BINARY)
	${BINARY}

test:
	go test -x -v zalua/dsl

tar: clean
	mkdir -p rpm/SOURCES
	tar --transform='s,^\.,$(NAME)-$(VERSION),'\
		-czf rpm/SOURCES/$(NAME)-$(VERSION).tar.gz .\
		--exclude=rpm/SOURCES

docker: submodule_check tar
	cp -a $(CURDIR)/rpm /build
	cp -a $(CURDIR)/rpm/SPECS/$(NAME).spec /build/SPECS/$(NAME)-$(VERSION).spec
	sed -i 's|%define version unknown|%define version $(VERSION)|g' /build/SPECS/$(NAME)-$(VERSION).spec
	chown -R root:root /build
	rpmbuild -ba --define '_topdir /build'\
		/build/SPECS/$(NAME)-$(VERSION).spec

clean:
	rm -f $(BINARY)
	rm -f rpm-tmp.*

.DEFAULT_GOAL: $(BINARY)

include Makefile.git
