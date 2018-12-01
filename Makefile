PKGNAME  = cuppa
DESTDIR ?=
PREFIX  ?= /usr
BINDIR   = $(PREFIX)/bin

GOBIN       = _build/bin
GOPROJROOT  = $(GOSRC)/$(PROJREPO)

GOLDFLAGS   = -ldflags "-s -w"
GOTAGS      = --tags "libsqlite3 linux"
GOCC        = go
GOFMT       = $(GOCC) fmt -x
GOGET       = $(GOCC) get $(GOLDFLAGS)
GOBUILD     = $(GOCC) build -v $(GOLDFLAGS) $(GOTAGS)
GOTEST      = $(GOCC) test
GOVET       = $(GOCC) vet
GOINSTALL   = $(GOCC) install $(GOLDFLAGS)
GOBUILDDEP  = GOPATH=`pwd`/_build $(GOINSTALL)
GOCLEANDEP  = GOPATH=`pwd`/_build $(GOCC) clean -cache -modcache

include Makefile.waterlog

GOLINT    = $(GOBIN)/golint -set_exit_status

all: build

build: setup-deps
	@$(call stage,BUILD)
	@$(GOBUILD)
	@$(call pass,BUILD)

test: build
	@$(call stage,TEST)
	@$(GOTEST) ./...
	@$(call pass,TEST)

validate: setup-deps
	@$(call stage,FORMAT)
	@$(GOFMT) ./...
	@$(call pass,FORMAT)
	@$(call stage,VET)
	@$(call task,Running 'go vet'...)
	@$(GOVET) ./...
	@$(call pass,VET)
	@$(call stage,LINT)
	@$(call task,Running 'golint'...)
	@$(GOLINT) `go list ./... | grep -v vendor`
	@$(call pass,LINT)

setup-deps:
	@$(call stage,DEPS)
	@if [ -d build/src/honnef.co ]; then rm -rf build/src/honnef.co; fi
	@if [ ! -e $(GOBIN)/golint ]; then \
	    $(call task,Installing golint...); \
	    $(GOBUILDDEP) github.com/golang/lint/golint; \
        $(GOCLEANDEP) ./...; \
	fi

install:
	@$(call stage,INSTALL)
	install -D -m 00755 $(PKGNAME) $(DESTDIR)$(BINDIR)/$(PKGNAME)
	@$(call pass,INSTALL)

uninstall:
	@$(call stage,UNINSTALL)
	rm -f $(DESTDIR)$(BINDIR)/$(PKGNAME)
	@$(call pass,UNINSTALL)

clean:
	@$(call stage,CLEAN)
	@$(call task,Removing _build directory...)
	@rm -rf _build
	@$(call task,Removing executable...)
	@rm $(PKGNAME)
	@$(call pass,CLEAN)
