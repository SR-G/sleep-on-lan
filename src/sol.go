package main

import (
	"github.com/integrii/flaggy"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var configuration = Configuration{}
var configurationFileName = "sol.json"
var configurationFileNameFromCommandLine string
var exit chan bool

func ExitDaemon() {
	Info.Println("Stopping daemon ...")
	exit <- true
}

func init() {
	// Init flag reader
	flaggy.String(&configurationFileNameFromCommandLine, "c", "config", "Configuration file to use (optional, default is 'sol.json' next to the binary)")
	flaggy.SetName("Sleep-On-LAN")
	flaggy.SetDescription("Daemon allowing to send a linux or windows computer to sleep")
	flaggy.DefaultParser.ShowHelpOnUnexpected = true
	flaggy.SetVersion(Version.Version())
	flaggy.Parse()

	// Init loggers
	PreInitLoggers()

	// Init channel allowing to exit when listeners can't be started
	exit = make(chan bool)
}

func main() {
	// Check which configuration file to use (either from --config, either sol.json file alongside the binary, either default values)
	var fullConfigurationFileName string
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))

	if configurationFileNameFromCommandLine != "" {
		if _, err := os.Stat(configurationFileNameFromCommandLine); err == nil {
			Info.Println("Will use configuration file provided through --config parameter, path is [" + configurationFileNameFromCommandLine + "]")
			fullConfigurationFileName = configurationFileNameFromCommandLine
		} else {
			Warning.Println("Configuration file provided through --config parameter not found on disk, path is [" + configurationFileNameFromCommandLine + "], will try default value")
			fullConfigurationFileName = dir + string(os.PathSeparator) + configurationFileName
		}
	} else {
		fullConfigurationFileName = dir + string(os.PathSeparator) + configurationFileName
	}

	// Loads configuration
	configuration.InitDefaultConfiguration()
	configuration.Load(fullConfigurationFileName)
	configuration.Parse()
	Info.Println("Application [" + Version.ApplicationName + "], version [" + Version.Version() + "]")

	// Display found IP/MAC
	Info.Println("Now starting sleep-on-lan, hardware IP/mac addresses are : ")
	for key, value := range LocalNetworkMap() {
		Info.Println(" - local IP adress [" + key + "], mac [" + value + "], reversed mac [" + ReverseMacAddress(value) + "]")
	}

	// Display commands found in configuration
	Info.Println("Available commands are : ")
	for _, command := range configuration.Commands {
		Info.Println(" - operation [" + command.Operation + "], command [" + command.Command + "], default [" + strconv.FormatBool(command.IsDefault) + "], type [" + command.CommandType + "]")
	}

	// Starts listeners, per configuration
	for _, listenerConfiguration := range configuration.listenersConfiguration {
		if listenerConfiguration.active {
			if strings.EqualFold(listenerConfiguration.nature, "UDP") {
				go ListenUDP(listenerConfiguration.port)
			} else if strings.EqualFold(listenerConfiguration.nature, "HTTP") {
				go ListenHTTP(listenerConfiguration.port)
			}
		}
	}

	// Blocks forever ... excepted if there are some start errors (depending on configuration)
	select {
	case <-exit:
		return
	}
}
