SOURCEDIR=.
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

BINARY=sol
PWD := $(shell pwd)

VERSION=1.1.0
VERSION_LABEL=RELEASE
PACKAGE=sleep-on-lan
DISTRIBUTION_PACKAGE=SleepOnLAN-${VERSION}-${VERSION_LABEL}
BUILD_TIME=$(shell date "+%FT%T%z")

LDFLAGS=-ldflags "-s -w -X 'main.BuildCommit=`git rev-parse HEAD`' -X 'main.BuildVersion=${VERSION}' -X 'main.BuildVersionLabel=${VERSION_LABEL}' -X 'main.BuildCompilationTimestamp=${BUILD_TIME}'"

ifeq ($(shell hostname),jupiter)
	DOCKER_IMAGE="tensin-app-golang"
else
	DOCKER_IMAGE="library/golang"
endif

.PHONY: install clean deploy run 

.ONESHELL: # Applies to every targets in the file!

# DEV / Quick build
build:
	cd src/
	go install ${LDFLAGS} sleep-on-lan

# DEV / Clean
clean:
	rm -rf bin

# DEV / Deploy configuration
conf:
	mkdir bin/
	cp resources/configuration/dev/sol-local-developement-configuration.json bin/

# DEV / Run binary with DEV configuration
run:
	setcap 'cap_net_bind_service=+ep' ${GOPATH}/bin/sleep-on-lan
	# ${GOPATH}/bin/sleep-on-lan -c resources/configuration/dev/sol-local-development-configuration.HS-hard-error.json
	${GOPATH}/bin/sleep-on-lan -c resources/configuration/dev/sol-local-development-configuration.json

# DEV / Create ZIP 
distribution: install
	mkdir -p bin/linux_amd64/ bin/linux_arm64/ bin/linux_arm/ bin/windows_amd64/ bin/windows_386/
	cp ${GOPATH}/bin/${PACKAGE} bin/linux_amd64/sol
	cp ${GOPATH}/bin/linux_arm64/${PACKAGE} bin/linux_arm64/sol
	cp ${GOPATH}/bin/linux_arm/${PACKAGE} bin/linux_arm/sol
	cp ${GOPATH}/bin/windows_amd64/${PACKAGE}.exe bin/windows_amd64/sol.exe
	cp ${GOPATH}/bin/windows_386/${PACKAGE}.exe bin/windows_386/sol.exe
	cp resources/configuration/default/sol-basic-configuration.json bin/linux_amd64/sol.json
	cp resources/configuration/default/sol-basic-configuration.json bin/linux_arm64/sol.json
	cp resources/configuration/default/sol-basic-configuration.json bin/linux_arm/sol.json
	cp resources/configuration/default/sol-basic-configuration.json bin/windows_amd64/sol.json
	cp resources/configuration/default/sol-basic-configuration.json bin/windows_386/sol.json
	cp resources/script/*.bat bin/windows_amd64
	cp resources/script/*.bat bin/windows_386
	cd bin/ 
	zip -r -9 ${DISTRIBUTION_PACKAGE}.zip ./linux_amd64/
	zip -r -9 ${DISTRIBUTION_PACKAGE}.zip ./linux_arm64/
	zip -r -9 ${DISTRIBUTION_PACKAGE}.zip ./linux_arm/
	zip -r -9 ${DISTRIBUTION_PACKAGE}.zip ./windows_amd64/
	zip -r -9 ${DISTRIBUTION_PACKAGE}.zip ./windows_386/

# DEV / Update GO modules
mod: 
	cd src/
	go mod tidy

# DEV / Format code
format:
	cd src/
	gofmt -w .

# DEV / Build all binaries
install: clean
	cd src/
	GOARCH=386 GOOS=windows go install ${LDFLAGS} ${PACKAGE}
	GOARCH=amd64 GOOS=windows go install ${LDFLAGS} ${PACKAGE}
	GOARCH=amd64 GOOS=linux go install ${LDFLAGS} -a -tags netgo -installsuffix netgo ${PACKAGE}
	GOARCH=arm64 GOOS=linux go install ${LDFLAGS} -a -tags netgo -installsuffix netgo ${PACKAGE}
	GOARCH=arm GOOS=linux go install ${LDFLAGS} -a -tags netgo -installsuffix netgo ${PACKAGE}

# DEV / Run docker
docker:
	docker run --rm -it --name "docker-sleep-on-lan-build" -p 8009:8009 -v ${PWD}:/go ${DOCKER_IMAGE} /bin/bash
