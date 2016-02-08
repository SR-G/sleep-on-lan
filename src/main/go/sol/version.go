package main

// Program has to be built with -ldflags "-X main.variable=VALUE"
// ex. : -i -ldflags "-X main.VERSION=1.0.0-SNAPSHOT -X main.BUILD_DATE=1970-01-01_00:00:01"
// ex. : -i -ldflags "-X main.VERSION=1.0.0-SNAPSHOT -X main.BUILD_DATE=\"%date\%"" (on windows, but doesn't work with LideIDE)
// ex. : -i -ldflags "-X main.VERSION=1.0.0-SNAPSHOT -X main.BUILD_DATE=$(date -u '+%Y-%m-%d_%H:%M:%S')" (on linux)

// @see http://stackoverflow.com/questions/11354518/golang-application-auto-build-versioning

// TODO(serge) : put a gradle or a cmake makefile around this

var APPLICATION_NAME string = "sleep-on-lan"
var VERSION string
var BUILD_DATE string
