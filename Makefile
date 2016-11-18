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

validate: golint-setup
    @printf "Formatting..."
    @go fmt ./...
    @printf "DONE\n"
    @printf "Vetting..."
    @go vet ./...
    @printf "DONE\n"
    @printf "Linting..."
    @${GOBIN}/golint -set_exit_status ./...
    @printf "DONE\n"

golint-setup:
    @if [ ! -e ${GOBIN}/golint ]; then \
        go get -u github.com/golang/lint/golint; \
        rm -rf ${GOPATH}/src/golang.org ${GOPATH}/src/github.com/golang ${GOPATH}/pkg; \
    fi

clean:
    @unlink ${PROJROOT}/cuppa
    @rm -r build
