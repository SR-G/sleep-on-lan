package main

import "os/user"

func (conf *Configuration) RegisterDefaultCommand() {
	defaultCommand := CommandConfiguration{Operation: "sleep", CommandType: COMMAND_TYPE_EXTERNAL, IsDefault: true, Command: "systemctl suspend"}
	conf.Commands = []CommandConfiguration{defaultCommand}
}

func RegisterPossibleConfigurationFileNames() []PossibleConfigurationFilename {
	var possibleConfigurationFileNames []PossibleConfigurationFilename
	possibleConfigurationFileNames = append(possibleConfigurationFileNames, PossibleConfigurationFilename{"/etc/sol.json", "default configuration filename under /etc/ (linux)"})
	possibleConfigurationFileNames = append(possibleConfigurationFileNames, PossibleConfigurationFilename{"/etc/sleep-on-lan.json", "default configuration filename under /etc/ (linux)"})
	return possibleConfigurationFileNames
}

func ExecuteCommand(Command CommandConfiguration) {
	if Command.CommandType == COMMAND_TYPE_EXTERNAL {
		usernameAsString := "UNDEFINED"
		user, _ := user.Current()
		if user != nil {
			usernameAsString = user.Name + "(" + user.Uid + ")"
		}
		logger.Infof("Executing operation [" + Command.Operation + "], type [" + Command.Command + "], command [" + Command.Command + "], current user [" + usernameAsString + "]")
		sleepCommandLineImplementation(Command.Command)
	} else {
		logger.Infof("Unknown command type [" + Command.CommandType + "]")
	}
}

func sleepCommandLineImplementation(cmd string) {
	if cmd == "" {
		cmd = "pm-suspend"
	}
	logger.Infof("Sleep implementation [linux], sleep command is [" + cmd + "]")
	_, output, err := Execute(cmd)
	if err != nil {
		logger.Errorf("Can't execute command ["+cmd+"], output is : "+output+", error is : ", err)
	} else {
		logger.Infof("Command correctly executed")
	}
}
