package main

import (
	"syscall"
)

func RegisterDefaultCommand() {
	defaultCommand := CommandConfiguration{Operation: "sleep", CommandType: "internal-dll", IsDefault: true}
	configuration.Commands = []CommandConfiguration{defaultCommand}
}

func ExecuteCommand(Command CommandConfiguration) {
	if Command.CommandType == "internal-dll" {
		Info.Println("Executing operation [" + Command.Operation + "], type[" + Command.CommandType + "]")
		if Command.Operation == "sleep" {
			sleepDLLImplementation()
		}
		if Command.Operation == "shutdown" {
			shutdownDLLImplementation()
		}
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

	Info.Printf("Command executed, result code [" + string(ret) + "]")
}

func shutdownDLLImplementation() {
	var mod = syscall.NewLazyDLL("Advapi32.dll")
	var proc = mod.NewProc("InitiateSystemShutdown")

	// DLL API : public static extern bool InitiateSystemShutdown(string lpMachineName, string lpMessage, int dwTimeout, bool bForceAppsClosed, bool bRebootAfterShutdown);
	// ex. : uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("Done Title"))),
	ret, _, _ := proc.Call(0,
		string(""), // lpMachineName
		string(""), // lpMessage
		uintptr(0), // dwTimeout
		uintptr(0), // bForceAppsClosed
		uintptr(0)) // bRebootAfterShutdown

	Info.Printf("Command executed, result code [" + string(ret) + "]")
}


