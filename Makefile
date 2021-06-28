SOURCEDIR=.
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

BINARY=sol
PWD := $(shell pwd)

VERSION=1.0.4-SNAPSHOT
PACKAGE=SleepOnLAN-${VERSION}
BUILD_TIME=$(date "%FT%T%z")

LDFLAGS=-ldflags "-d -s -w -X tensin.org/sol/core/version.Build=`git rev-parse HEAD`" -a -tags netgo -installsuffix netgo
ifeq ($(shell hostname),jupiter)
	DOCKER_IMAGE="tensin-app-golang"
else
	DOCKER_IMAGE="library/golang"
endif

.PHONY: install clean deploy run 

build:
	go install sleep-on-lan

clean:
	rm -rf bin

conf:
	cp resources/sol.json bin/

run:
	bin/sol

distribution: install
	mkdir bin/linux/ 
	mv bin/sol bin/linux
	cp resources/sol.json bin/linux/ 
	cp resources/sol.json bin/windows_amd64/
	cp resources/script/*.bat bin/windows_amd64
	cd bin/ ; zip -r -9 ${PACKAGE}.zip ./linux ; zip -r -9 ${PACKAGE}.zip ./windows_amd64

install: clean
	rm -rf bin
	GOARCH=amd64 GOOS=windows go install sleep-on-lan
	GOARCH=amd64 GOOS=linux go install -ldflags "-d -s -w -X tensin.org/watchthatpage/core.Build=`git rev-parse HEAD`" -a -tags netgo -installsuffix netgo sleep-on-lan

docker:
	docker run --rm -it -v ${PWD}:/go ${DOCKER_IMAGE} /bin/bash
