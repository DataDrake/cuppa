.RECIPEPREFIX != ps

GOPATH   = $(shell pwd)/build
GOBIN    = ${GOPATH}/bin
GOSRC    = ${GOPATH}/src
PROJROOT = ${GOSRC}/github.com/DataDrake

all: build

build: setup
    @printf "Building..."
    @go install github.com/DataDrake/cuppa
    @printf "DONE\n"

setup:
    @printf "Setting up..."
    @if [ ! -d ${GOPATH} ]; then mkdir -p $(GOPATH); fi
    @if [ ! -d ${GOSRC} ]; then mkdir -p ${GOSRC}; fi
    @if [ ! -d ${PROJROOT} ]; then mkdir -p ${PROJROOT}; fi
    @if [ ! -d ${PROJROOT}/cuppa ]; then ln -s $(shell pwd) ${PROJROOT}/cuppa; fi
    @printf "DONE\n"

validate:
    go fmt ./...
    go vet ./...
    @if [ ! -d ${GOPATH}/bin/golint ]; then go get -u github.com/golang/lint/golint; fi
    ${GOBIN}/golint github.com/DataDrake/cuppa

clean:
    @unlink ${PROJROOT}/cuppa
    @rm -r build
