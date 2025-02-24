package main

import (
	"bytes"
	"fmt"
)

// These variables are injected through Makefile at compile time
var BuildVersion = "1.1.2"
var BuildVersionLabel = "SNAPSHOT"
var BuildCommit = ""
var BuildCompilationTimestamp = ""

type version struct {
	ApplicationName      string // Name of the program
	Version              string // Version, e.g. 1.0.0
	VersionLabel         string // Label of the version, like SNAPSHOT or RELEASE
	VersionName          string // Name of the version, like Buster, ...
	Commit               string // GiT commit hash
	CompilationTimestamp string // When program has been compiled
}

// Version string
var Version = version{ApplicationName: "sleep-on-lan", Version: BuildVersion, VersionLabel: BuildVersionLabel, Commit: BuildCommit, CompilationTimestamp: BuildCompilationTimestamp}

func (v version) DumpVersion() string {
	var buf bytes.Buffer
	buf.WriteString(v.Version)
	if v.VersionLabel != "" {
		buf.WriteString("-" + v.VersionLabel)
	}
	return buf.String()
}

func (v version) GetVersion() string {
	var buf bytes.Buffer
	buf.WriteString(v.Version)
	if v.VersionLabel != "" {
		buf.WriteString("-" + v.VersionLabel)
	}
	return buf.String()
}

func (v version) String() string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%s %s", v.ApplicationName, v.Version))
	if v.VersionLabel != "" {
		buf.WriteString("-" + v.VersionLabel)
	}
	if v.VersionName != "" {
		buf.WriteString(" \"" + v.VersionName + "\"")
	}
	if v.CompilationTimestamp != "" {
		buf.WriteString(" (" + v.CompilationTimestamp + ")")
	}
	if v.Commit != "" {
		buf.WriteString("\nGit commit hash: " + v.Commit)
	}
	return buf.String()
}
