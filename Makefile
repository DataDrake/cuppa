.RECIPEPREFIX != ps

GOCC     = go

GOPATH   = $(shell pwd)/build
GOBIN    = build/bin
GOSRC    = build/src
PROJROOT = $(GOSRC)/github.com/DataDrake

DESTDIR ?=
PREFIX  ?= /usr
BINDIR   = $(PREFIX)/bin

all: build

build: setup
    @printf "\e[34m[Start Build]\e[39m\n"
    $(GOCC) install -v github.com/DataDrake/cuppa
    @printf "\e[34m[End Build]\e[39m\n\n"

setup:
    @printf "\e[34m[Start Setup]\e[39m\n"
    @printf "Setting up GOPATH...\n"
    @mkdir -p $(GOPATH)
    @printf "Setting up src/...\n"
    @mkdir -p $(GOSRC)
    @printf "Setting up project root...\n"
    @mkdir -p $(PROJROOT)
    @printf "Setting up symlinks...\n"
    @if [ ! -d $(PROJROOT)/cuppa ]; then ln -s $(shell pwd) $(PROJROOT)/cuppa; fi
    @printf "\e[34m[End Setup]\e[39m\n\n"

validate: golint-setup
    @printf "\e[34m[Start Format]\e[39m\n"
    $(GOCC) fmt -x ./...
    @printf "\e[34m[End Format]\e[39m\n\n"
    @printf "\e[34m[Start Vet]\e[39m\n"
    $(GOCC) vet -x ./...
    @printf "\e[34m[End Vet]\e[39m\n\n"
    @printf "\e[34m[Start Lint]\e[39m\n"
    $(GOBIN)/golint -set_exit_status ./...
    @printf "\e[34m[End Lint]\e[39m\n\n"

golint-setup:
    @if [ ! -e $(GOBIN)/golint ]; then \
        printf "Installing golint..."; \
        $(GOCC) get -u github.com/golang/lint/golint; \
        printf "DONE\n\n"; \
        rm -rf $(GOPATH)/src/golang.org $(GOPATH)/src/github.com/golang $(GOPATH)/pkg; \
    fi

install:
    @printf "\e[34m[Start Install]\e[39m\n"
    install -D -m 00755 $(GOBIN)/cuppa $(DESTDIR)$(BINDIR)/cuppa
    @printf "\e[34m[End Install]\e[39m\n"

uninstall:
    @printf "\e[34m[Start Uninstall]\e[39m\n"
    rm -f $(DESTDIR)$(BINDIR)/cuppa
    @printf "\e[34m[End Uninstall]\e[39m\n"

clean:
    @printf "\e[34m[Start Clean]\e[39m\n"
    @printf "Removing symlinks...\n"
    unlink $(PROJROOT)/cuppa
    @printf "Removing build directory...\n"
    rm -r build
    @printf "\e[34m[End Clean]\e[39m\n\n"
