SOURCEDIR=.
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

BINARY=sol
PWD := $(shell pwd)

VERSION=1.0.0
BUILD_TIME=$(date "%FT%T%z")

LDFLAGS=-ldflags "-d -s -w -X tensin.org/sol/core/version.Build=`git rev-parse HEAD`" -a -tags netgo -installsuffix netgo
PACKAGE=tensin.org/sol

$(BINARY): $(SOURCES)
	        go build ${LDFLAGS} -o bin/${BINARY} ${PACKAGE}

.PHONY: install clean deploy run 

build:
	        time go install ${PACKAGE}

install:
	        time go install ${LDFLAGS} ${PACKAGE}

deploy:
	        cp bin/webscrapper /home/bin/

clean:
	        rm bin/cache/*
	        [ -f bin/${BINARY} ] && rm -f bin/${BINARY}

run:
	        bin/sol

test:
	        go test -v tensin.org/webscrapper/core

docker:
	        docker run --rm -it -v ${PWD}:/go tensin-app-golang /bin/bash
