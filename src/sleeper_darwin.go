package main

func (conf *Configuration) RegisterDefaultCommand() {
	defaultCommand := CommandConfiguration{Operation: "sleep", CommandType: COMMAND_TYPE_EXTERNAL, IsDefault: true, Command: "pmset sleepnow"}
	conf.Commands = []CommandConfiguration{defaultCommand}
}

func RegisterPossibleConfigurationFileNames() []PossibleConfigurationFilename {
	var possibleConfigurationFileNames []PossibleConfigurationFilename
	possibleConfigurationFileNames = append(possibleConfigurationFileNames, PossibleConfigurationFilename{"/etc/sol.json", "default configuration filename under /etc/ (darwin)"})
	possibleConfigurationFileNames = append(possibleConfigurationFileNames, PossibleConfigurationFilename{"/etc/sleep-on-lan.json", "default configuration filename under /etc/ (darwin)"})
	return possibleConfigurationFileNames
}

func ExecuteCommand(Command CommandConfiguration) {
	if Command.CommandType == COMMAND_TYPE_EXTERNAL {
		logger.Infof("Executing operation [" + Command.Operation + "], type [" + Command.Command + "], command [" + Command.Command + "]")
		sleepCommandLineImplementation(Command.Command)
	} else {
		logger.Infof("Unknown command type [" + Command.CommandType + "]")
	}
}

func sleepCommandLineImplementation(cmd string) {
	if cmd == "" {
		cmd = "pmset sleepnow"
	}
	logger.Infof("Sleep implementation [darwin], sleep command is [" + cmd + "]")
	_, _, err := Execute(cmd)
	if err != nil {
		logger.Errorf("Can't execute command [" + cmd + "] : " + err.Error())
	} else {
		logger.Infof("Command correctly executed")
	}
}
