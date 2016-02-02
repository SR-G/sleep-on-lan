package main

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var configuration = Configuration{}
var configurationFileName = "sol.json"

func main() {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	PreInitLoggers()
	var fullConfigurationFileName = dir + string(os.PathSeparator) + configurationFileName
	configuration.InitDefaultConfiguration()
	configuration.Load(fullConfigurationFileName)
	configuration.Parse()
	Info.Println("Loaded configuration from [" + fullConfigurationFileName + "]")

	Info.Println("Now starting sleep-on-lan, hardware IP/mac addresses are : ")
	for key, value := range LocalNetworkMap() {
		Info.Println(" - local IP adress [" + key + "], mac [" + value + "]")
	}

	for _, command := range configuration.Commands {
		Info.Println("  - operation [" + command.Operation + "], command [" + command.Command + "], default [" + strconv.FormatBool(command.IsDefault) + "], type [" + command.CommandType + "]")
	}

	for _, listenerConfiguration := range configuration.listenersConfiguration {
		if listenerConfiguration.active {
			if strings.EqualFold(listenerConfiguration.nature, "UDP") {
				go ListenUDP(listenerConfiguration.port)
			} else if strings.EqualFold(listenerConfiguration.nature, "HTTP") {
				go ListenHTTP(listenerConfiguration.port, configuration.Commands)
			}
		}
	}

	// Info.Println("sleep-on-lan up and running")
	select {} // block forever
}
