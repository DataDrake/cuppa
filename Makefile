PKGNAME  = cuppa
DESTDIR ?=
PREFIX  ?= /usr
BINDIR   = $(PREFIX)/bin

GOBIN       = ~/go/bin
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
	@$(GOLINT)
	@$(call pass,LINT)

setup-deps:
	@$(call stage,DEPS)
	@if [ ! -e $(GOBIN)/golint ]; then \
	    $(call task,Installing golint...); \
	    $(GOGET) golang.org/x/lint; \
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
	@$(call task,Removing executable...)
	@rm $(PKGNAME)
	@$(call pass,CLEAN)
