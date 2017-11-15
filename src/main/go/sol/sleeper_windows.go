package main

import (
	"syscall"
	"fmt"
)

func RegisterDefaultCommand() {
	defaultCommand := CommandConfiguration{Operation: "sleep", CommandType: "internal-dll", IsDefault: true}
	configuration.Commands = []CommandConfiguration{defaultCommand}
}

func ExecuteCommand(Command CommandConfiguration) {
	if Command.CommandType == "internal-dll" {
		Info.Println("Executing operation [" + Command.Operation + "], type[" + Command.CommandType + "]")
		sleepDLLImplementation()
	} else {
		Info.Println("Executing operation [" + Command.Operation + "], type[" + Command.CommandType + "], command [" + Command.Command + "]")
		sleepCommandLineImplementation(Command.Command)
	}
}

func sleepCommandLineImplementation(cmd string) {
	if cmd == "" {
		cmd = "C:\\Windows\\System32\\rundll32.exe powrprof.dll,SetSuspendState 0,1,1"
	}
	Info.Println("Sleep implementation [windows], sleep command is [", cmd, "]")
	_, _, err := Execute(cmd)
	if err != nil {
		Error.Println("Can't execute command [" + cmd + "] : " + err.Error())
	} else {
		Info.Println("Command correctly executed")
	}
}

func sleepDLLImplementation() {
	var mod = syscall.NewLazyDLL("Powrprof.dll")
	var proc = mod.NewProc("SetSuspendState")

	// DLL API : public static extern bool SetSuspendState(bool hiberate, bool forceCritical, bool disableWakeEvent);
	// ex. : uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("Done Title"))),
	ret, _, _ := proc.Call(0,
		uintptr(0), // hibernate
		uintptr(1), // forceCritical
		uintptr(1)) // disableWakeEvent
	Info.Printf("Command executed, result code [" + fmt.Sprint(ret) + "]")
}
