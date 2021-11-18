package main

import (
	"bytes"
	"fmt"
)

// These variables are injected through Makefile at compile time
var BuildVersion = "1.0.7"
var BuildVersionLabel = "SNAPSHOT"
var BuildCommit = ""

type version struct {
	ApplicationName string // Name of the program
	Version         string // Version, e.g. 1.0.0
	VersionLabel    string // Label of the version, like SNAPSHOT or RELEASE
	VersionName     string // Name of the version, like Buster, ...
	Build           string // GiT commit hash
}

// Version string
var Version = version{ApplicationName: "sleep-on-lan", Version: BuildVersion, VersionLabel: BuildVersionLabel, Build: BuildCommit}

func (v version) DumpVersion() string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%s", v.Version))
	if v.VersionLabel != "" {
		buf.WriteString("-" + v.VersionLabel)
	}
	return buf.String()
}

func (v version) String() string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%s version %s", v.ApplicationName, v.Version))
	if v.VersionLabel != "" {
		buf.WriteString("-" + v.VersionLabel)
	}
	if v.VersionName != "" {
		buf.WriteString(" \"" + v.VersionName + "\"")
	}
	if v.Build != "" {
		buf.WriteString("\nGit commit hash: " + v.Build)
	}
	return buf.String()
}
