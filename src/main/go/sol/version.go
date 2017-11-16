package main

import (
	"bytes"
	"fmt"
)

type version struct {
	ApplicationName     string
	Major, Minor, Patch int
	VersionLabel        string
	VersionName         string
}

// Version string
var Version = version{"sleep-on-lan", 1, 0, 2, "SNAPSHOT", ""}

// Build string
var Build string

func (v version) Version() string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch))
	if v.VersionLabel != "" {
		buf.WriteString("-" + v.VersionLabel)
	}
	return buf.String()
}

func (v version) String() string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%s version %d.%d.%d", v.ApplicationName, v.Major, v.Minor, v.Patch))
	if v.VersionLabel != "" {
		buf.WriteString("-" + v.VersionLabel)
	}
	if v.VersionName != "" {
		buf.WriteString(" \"" + v.VersionName + "\"")
	}
	if Build != "" {
		buf.WriteString("\nGit commit hash: " + Build)
	}
	return buf.String()
}
