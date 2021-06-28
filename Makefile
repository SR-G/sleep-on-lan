SOURCEDIR=.
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

BINARY=sol
PWD := $(shell pwd)

VERSION=1.0.5-SNAPSHOT
PACKAGE=SleepOnLAN-${VERSION}
BUILD_TIME=$(date "%FT%T%z")

LDFLAGS=-ldflags "-d -s -w -X sleep-on-lan/version.Build=`git rev-parse HEAD`" -a -tags netgo -installsuffix netgo
ifeq ($(shell hostname),jupiter)
	DOCKER_IMAGE="tensin-app-golang"
else
	DOCKER_IMAGE="library/golang"
endif

.PHONY: install clean deploy run 

.ONESHELL: # Applies to every targets in the file!

build:
	cd src/
	go install sleep-on-lan

clean:
	rm -rf bin

conf:
	mkdir bin/
	cp resources/sol.json bin/

run:
	${GOPATH}/bin/sleep-on-lan

distribution: install
	mkdir -p bin/linux/ bin/windows_amd64/
	cp ${GOPATH}/bin/sleep-on-lan bin/linux/sol
	cp ${GOPATH}/bin/windows_amd64/sleep-on-lan.exe bin/windows_amd64/sol.exe
	cp resources/sol.json bin/linux/ 
	cp resources/sol.json bin/windows_amd64/
	cp resources/script/*.bat bin/windows_amd64
	cd bin/ ; zip -r -9 ${PACKAGE}.zip ./linux ; zip -r -9 ${PACKAGE}.zip ./windows_amd64

install: clean
	cd src/
	GOARCH=amd64 GOOS=windows go install sleep-on-lan
	GOARCH=amd64 GOOS=linux go install -ldflags "-d -s -w -X tensin.org/watchthatpage/core.Build=`git rev-parse HEAD`" -a -tags netgo -installsuffix netgo sleep-on-lan

docker:
	docker run --rm -it -v ${PWD}:/go ${DOCKER_IMAGE} /bin/bash
