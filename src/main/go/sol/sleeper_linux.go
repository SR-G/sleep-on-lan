package main

import (
	// "fmt"
	// "os/exec"
)

func Sleep() {
	Info.Println("Sleep linux")
	sleepCommandLineImplementation()
}

func sleepCommandLineImplementation() {
	var cmd = ""
	if configuration.SleepCommand != "" {
		cmd = configuration.SleepCommand
	} else {
		cmd = "pm-suspend"
	}
	Info.Println("Sleep implementation [linux], sleep command is [" + cmd + "]")
	_, _, err := Execute(cmd)
	if err != nil {
		Error.Println("Can't execute command [" + cmd + "] : " + err.Error())
	} else {
		Info.Println("Command correctly executed")
	}
}
