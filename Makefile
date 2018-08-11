PKGNAME  = cuppa
SUBPKGS  = cli config providers results
PROJREPO = github.com/DataDrake

include Makefile.golang
include Makefile.waterlog

MEGACHECK = $(GOBIN)/megacheck
GOLINT    = $(GOBIN)/golint -set_exit_status

DESTDIR ?=
PREFIX  ?= /usr
BINDIR   = $(PREFIX)/bin

all: build

build: setup setup-deps
	@$(call stage,BUILD)
	@$(GOINSTALL) $(PROJREPO)/$(PKGNAME)
	@$(call pass,BUILD)

setup:
	@$(call stage,SETUP)
	@$(call task,Setting up project root...)
	@mkdir -p $(GOPROJROOT)
	@$(call task,Setting up symlinks...)
	@if [ ! -d $(GOPROJROOT)/$(PKGNAME) ]; then ln -s $(shell pwd) $(GOPROJROOT)/$(PKGNAME); fi
	@$(call pass,SETUP)

test: build
	@$(call stage,TEST)
	@for d in $(SUBPKGS); do $(GOTEST) ./$$d/... || exit 1; done
	@$(call pass,TEST)

validate: setup-deps
	@$(call stage,FORMAT)
	@for d in $(SUBPKGS); do $(GOFMT) ./$$d/...|| exit 1; done || $(GOFMT) $(PKGNAME).go
	@$(call pass,FORMAT)
	@$(call stage,VET)
	@$(call task,Running 'go vet'...)
	@cd $(GOPROJROOT)/$(PKGNAME); for d in $(SUBPKGS); do $(GOVET) ./... && exit 1; done || $(GOVET) $(PKGNAME).go || exit 1
	@$(call task,Running 'megacheck'...)
	@for d in $(SUBPKGS); do $(MEGACHECK) ./$$d || exit 1; done || $(MEGACHECK) $(PKGNAME).go || exit 1
	@$(call pass,VET)
	@$(call stage,LINT)
	@$(call task,Running 'golint'...)
	@for d in $(SUBPKGS); do $(GOLINT) ./$$d/... || exit 1; done || $(GOLINT) $(PKGNAME).go || exit 1;
	@$(call pass,LINT)

setup-deps:
	@$(call stage,DEPS)
	@if [ ! -e $(GOBIN)/dep ]; then \
	    $(call task,Installing dep...); \
	    $(GOGET) -d github.com/golang/dep/cmd/dep; \
	    cd build/src/github.com/golang/dep/cmd/dep; \
	    git checkout tags/v0.4.1; \
	    $(GOINSTALL) ./...; \
	    cd $(GOPROJROOT)/$(PKGNAME); \
	fi
	@if [ ! -e $(GOBIN)/megacheck ]; then \
	    $(call task,Installing megacheck...); \
	    $(GOGET) honnef.co/go/tools/cmd/megacheck; \
	fi
	@if [ -d build/src/honnef.co ]; then rm -rf build/src/honnef.co; fi
	@if [ ! -e $(GOBIN)/golint ]; then \
	    $(call task,Installing golint...); \
	    $(GOGET) github.com/golang/lint/golint; \
	fi
	@if [ -d build/src/golang.org ]; then rm -rf build/src/golang.org; fi
	@if [ -d build/src/github.com/golang ]; then rm -rf build/src/github.com/golang; fi
	@if [ ! -d vendor ]; then \
	    $(call task,Getting build dependencies...); \
	    cd $(GOPROJROOT)/$(PKGNAME); GOPATH=$(GOPATH) $(GOBIN)/dep ensure; \
	fi

install:
	@$(call stage,INSTALL)
	install -D -m 00755 $(GOBIN)/$(PKGNAME) $(DESTDIR)$(BINDIR)/$(PKGNAME)
	@$(call pass,INSTALL)

uninstall:
	@$(call stage,UNINSTALL)
	rm -f $(DESTDIR)$(BINDIR)/$(PKGNAME)
	@$(call pass,UNINSTALL)

clean:
	@$(call stage,CLEAN)
	@$(call task,Removing symlinks...)
	@unlink $(GOPROJROOT)/$(PKGNAME)
	@$(call task,Removing build directory...)
	@rm -rf build
	@$(call pass,CLEAN)
