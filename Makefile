SOURCEDIR=.
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

BINARY=sol
PWD := $(shell pwd)

VERSION=1.0.2-SNAPSHOT
PACKAGE=SleepOnLAN-${VERSION}
BUILD_TIME=$(date "%FT%T%z")

LDFLAGS=-ldflags "-d -s -w -X tensin.org/sol/core/version.Build=`git rev-parse HEAD`" -a -tags netgo -installsuffix netgo

.PHONY: install clean deploy run 

build:
	        cd /go/src/
			go install main/go/sol/

clean:
	        rm -rf bin

run:
	        bin/sol

distribution: install
			mkdir /go/bin/linux/ 
			mv /go/bin/sol /go/bin/linux
			cp /go/src/main/resources/sol.json /go/bin/linux/ 
			cp /go/src/main/resources/sol.json /go/bin/windows_amd64/
			cp /go/src/script/*.bat /go/bin/windows_amd64
			cd /go/bin/ ; zip -r -9 ${PACKAGE}.zip ./linux ; zip -r -9 ${PACKAGE}.zip ./windows_amd64

install: clean
			rm -rf /go/bin
			cd /go/src
			GOARCH=amd64 GOOS=windows go install main/go/sol/
			GOARCH=amd64 GOOS=linux go install -ldflags "-d -s -w -X tensin.org/watchthatpage/core.Build=`git rev-parse HEAD`" -a -tags netgo -installsuffix netgo main/go/sol/

docker:
	        docker run --rm -it -v ${PWD}:/go tensin-app-golang /bin/bash
