include Makefile.waterlog

GOPATH   = $(shell pwd)/build
GOCC     = GOPATH=$(GOPATH) go

GOBIN    = build/bin
GOSRC    = build/src
PROJROOT = $(GOSRC)/github.com/DataDrake
PKGNAME  = cuppa
SUBPKGS  = cmd \
           providers \
           results

DEPS     = github.com/DataDrake/cli-ng \
           github.com/DataDrake/waterlog

DESTDIR ?=
PREFIX  ?= /usr
BINDIR   = $(PREFIX)/bin

all: build

build: setup
	@$(call stage,BUILD)
	@$(GOCC) install -v -ldflags '-s -w' github.com/DataDrake/$(PKGNAME)
	@$(call pass,BUILD)

setup:
	@$(call stage,SETUP)
	@$(call task,Setting up GOPATH...)
	@mkdir -p $(GOPATH)
	@$(call task,Setting up src/...)
	@mkdir -p $(GOSRC)
	@$(call task,Setting up project root...)
	@mkdir -p $(PROJROOT)
	@$(call task,Setting up symlinks...)
	@if [ ! -d $(PROJROOT)/$(PKGNAME) ]; then ln -s $(shell pwd) $(PROJROOT)/$(PKGNAME); fi
	@$(call task,Getting dependencies...)
	@for d in $(DEPS); do $(GOCC) get $$d || exit 1; done
	@$(call pass,SETUP)

test: build
	@$(call stage,TEST)
	@for d in $(SUBPKGS); do $(GOCC) test -cover ./$$d/... || exit 1; done
	@$(call pass,TEST)

validate: golint-setup
	@$(call stage,FORMAT)
	@for d in $(SUBPKGS); do $(GOCC) fmt -x ./$$d/...|| exit 1; done || $(GOCC) fmt -x $(PKGNAME).go
	@$(call pass,FORMAT)
	@$(call stage,VET)
	@for d in $(SUBPKGS); do $(GOCC) vet -x ./$$d/...|| exit 1; done || $(GOCC) vet -x $(PKGNAME).go
	@$(call pass,VET)
	@$(call stage,LINT)
	@for d in $(SUBPKGS); do $(GOBIN)/golint -set_exit_status ./$$d/... || exit 1; done || $(GOBIN)/golint -set_exit_status $(PKGNAME).go || exit 1;
	@$(call pass,LINT)

golint-setup:
	@if [ ! -e $(GOBIN)/golint ]; then \
	    printf "Installing golint..."; \
	    $(GOCC) get -u github.com/golang/lint/golint; \
	    printf "DONE\n\n"; \
	    rm -rf $(GOPATH)/src/golang.org $(GOPATH)/src/github.com/golang $(GOPATH)/pkg; \
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
	@unlink $(PROJROOT)/$(PKGNAME)
	@$(call task,Removing build directory...)
	@rm -rf build
	@$(call pass,CLEAN)
