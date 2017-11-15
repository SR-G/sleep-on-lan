package main

import (
	"bytes"
	"fmt"
)

type version struct {
	Major, Minor, Patch int
	Label               string
	Name                string
}

const projectName = "sleep-on-lan"

// Version string
var Version = version{1, 0, 0, "SNAPSHOT", "FIRST_ITERATION"}

// Build string
var Build string

func (v version) String() string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%s version %d.%d.%d", projectName, v.Major, v.Minor, v.Patch))
	if v.Label != "" {
		buf.WriteString("-" + v.Label)
	}
	if v.Name != "" {
		buf.WriteString(" \"" + v.Name + "\"")
	}
	if Build != "" {
		buf.WriteString("\nGit commit hash: " + Build)
	}
	return buf.String()
}
