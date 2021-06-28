package main

func RegisterDefaultCommand() {
	defaultCommand := CommandConfiguration{Operation: "sleep", CommandType: COMMAND_TYPE_EXTERNAL, IsDefault: true, Command: "pm-suspend"}
	configuration.Commands = []CommandConfiguration{defaultCommand}
}

func ExecuteCommand(Command CommandConfiguration) {
	if Command.CommandType == COMMAND_TYPE_EXTERNAL {
		Info.Println("Executing operation [" + Command.Operation + "], type [" + Command.Command + "], command [" + Command.Command + "]")
		sleepCommandLineImplementation(Command.Command)
	} else {
		Info.Println("Unknown command type [" + Command.CommandType + "]")
	}
}

func sleepCommandLineImplementation(cmd string) {
	if cmd == "" {
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
