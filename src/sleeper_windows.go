package main

import (
	"fmt"
	"syscall"
	"unsafe"

	winio "github.com/Microsoft/go-winio"
)

const (
	DEFAULT_COMMAND_SLEEP    = "sleep"
	DEFAULT_COMMAND_SHUTDOWN = "shutdown"
)

func RegisterDefaultCommand() {
	defaultSleepCommand := CommandConfiguration{Operation: DEFAULT_COMMAND_SLEEP, CommandType: COMMAND_TYPE_INTERNAL_DLL, IsDefault: true}
	defaultShutdownCommand := CommandConfiguration{Operation: DEFAULT_COMMAND_SHUTDOWN, CommandType: COMMAND_TYPE_INTERNAL_DLL, IsDefault: false}
	configuration.Commands = []CommandConfiguration{defaultSleepCommand, defaultShutdownCommand}
}

func ExecuteCommand(Command CommandConfiguration) {
	if Command.CommandType == COMMAND_TYPE_INTERNAL_DLL {
		Info.Println("Executing operation [" + Command.Operation + "], type[" + Command.CommandType + "]")
		if Command.Operation == DEFAULT_COMMAND_SLEEP {
			sleepDLLImplementation()
		} else if Command.Operation == DEFAULT_COMMAND_SHUTDOWN {
			shutdownDLLImplementation()
		}
	} else if Command.CommandType == COMMAND_TYPE_EXTERNAL {
		Info.Println("Executing operation [" + Command.Operation + "], type[" + Command.CommandType + "], command [" + Command.Command + "]")
		sleepCommandLineImplementation(Command.Command)
	} else {
		Info.Println("Unknown command type [" + Command.CommandType + "]")
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

func shutdownDLLImplementation() {
	// SeShutdownPrivilege
	err := winio.RunWithPrivilege("SeShutdownPrivilege", func() error {
		var mod = syscall.NewLazyDLL("Advapi32.dll")
		var proc = mod.NewProc("InitiateSystemShutdownW")

		// DLL API : public static extern bool InitiateSystemShutdown(string lpMachineName, string lpMessage, int dwTimeout, bool bForceAppsClosed, bool bRebootAfterShutdown);
		// ex. : uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("Done Title"))),

		// var a [1]byte
		// a[0] = byte(0)
		// addrPtr := unsafe.Pointer(&a)
		ret, _, _ := proc.Call(0,
			uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(""))), // lpMachineName
			uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(""))), // lpMessage
			uintptr(0), // dwTimeout
			uintptr(1), // bForceAppsClosed
			uintptr(0)) // bRebootAfterShutdown

		// ret 0 = false, ret 1 = true = success
		Info.Printf("Command executed, result code [" + fmt.Sprint(ret) + "]")
		return nil
	})
	if err != nil {
		Error.Printf("Can't execute command")
	}
}
