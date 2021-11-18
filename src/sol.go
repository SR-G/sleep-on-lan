package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/integrii/flaggy"
	"github.com/labstack/gommon/color"

	"github.com/juju/loggo"
)

var configuration = Configuration{}
var configurationFileName = "sol.json"
var configurationFileNameFromCommandLine string
var subCommandGenerateConfiguration *flaggy.Subcommand
var exit chan bool
var colorer *color.Color
var logger = loggo.GetLogger("")

func ExitDaemon() {
	logger.Infof("Stopping daemon ...")
	exit <- true
}

func init() {

	// Preinit loggers with default value (will be later overridden per configuration, if needed)
	logger.SetLogLevel(loggo.INFO)

	// Init flag reader
	flaggy.String(&configurationFileNameFromCommandLine, "c", "config", "Configuration file to use (optional, default is 'sol.json' next to the binary)")
	flaggy.SetName("Sleep-On-LAN")
	flaggy.SetDescription("Daemon allowing to send a linux or windows computer to sleep")
	flaggy.DefaultParser.ShowHelpOnUnexpected = true
	flaggy.SetVersion(Version.String())

	subCommandGenerateConfiguration = flaggy.NewSubcommand("generate-configuration")
	subCommandGenerateConfiguration.Description = "Generate a default configuration JSON file"
	flaggy.AttachSubcommand(subCommandGenerateConfiguration, 1)

	flaggy.Parse()

	// Colors ...
	colorer = color.New()

	// Init channel allowing to exit when listeners can't be started
	exit = make(chan bool)
}

// Sub-command : generate JSON blank configuration (init)
func executeCommandGenerateConfiguration() {
	c := Configuration{}
	c.InitDefaultConfiguration()
	c.Parse()

	if _, err := os.Stat(configurationFileName); err == nil {
		logger.Errorf("Can't generate default JSON configuration info [" + configurationFileName + "] : file already exist")
	} else {
		logger.Infof("Writing default JSON configuration into [" + configurationFileName + "]")

		file, err := os.OpenFile(configurationFileName, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
		if err != nil {
			logger.Errorf("Can't write file ["+configurationFileName+"]", err)
		}
		defer file.Close()

		// Create blank configuration (with default values) and write bytes to file
		b, err := json.MarshalIndent(c, "", "    ")
		if err != nil {
			logger.Errorf("Can't create JSON content", err)
		}
		bytesWritten, err := file.Write(b)
		if err != nil {
			logger.Errorf("Can't write JSON content into ["+configurationFileName+"]", err)
		} else {
			logger.Infof("Bytes written : " + strconv.Itoa(bytesWritten))
		}
	}
}

func startDaemon() {
	// Check which configuration file to use (either from --config, either sol.json file alongside the binary, either default values)
	var fullConfigurationFileName string
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))

	if configurationFileNameFromCommandLine != "" {
		if _, err := os.Stat(configurationFileNameFromCommandLine); err == nil {
			logger.Infof("Will use configuration file provided through --config parameter, path is [" + colorer.Green(configurationFileNameFromCommandLine) + "]")
			fullConfigurationFileName = configurationFileNameFromCommandLine
		} else {
			logger.Warningf("Configuration file provided through --config parameter not found on disk, path is [" + colorer.Red(configurationFileNameFromCommandLine) + "], will try default value")
			fullConfigurationFileName = dir + string(os.PathSeparator) + configurationFileName
		}
	} else {
		fullConfigurationFileName = dir + string(os.PathSeparator) + configurationFileName
	}

	// Loads configuration
	configuration.InitDefaultConfiguration()
	configuration.Load(fullConfigurationFileName)
	configuration.RefineLogger()
	configuration.Parse()
	logger.Infof("Application [" + colorer.Green(Version.ApplicationName) + "], version [" + colorer.Green(Version.String()) + "]")

	// Display found IP/MAC
	logger.Infof("Now starting sleep-on-lan, hardware IP/mac addresses are : ")
	for key, value := range LocalNetworkMap() {
		logger.Infof(" - local IP adress [" + colorer.Green(key) + "], mac [" + colorer.Green(value) + "], reversed mac [" + colorer.Green(ReverseMacAddress(value)) + "]")
	}

	// Display commands found in configuration
	logger.Infof("Available commands are : ")
	for _, command := range configuration.Commands {
		logger.Infof(" - operation [" + color.Green(command.Operation) + "], command [" + color.Green(command.Command) + "], default [" + color.Green(strconv.FormatBool(command.IsDefault)) + "], type [" + color.Green(command.CommandType) + "]")
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

func main() {
	// Handling subcommands, if any
	if subCommandGenerateConfiguration.Used {
		executeCommandGenerateConfiguration()
	} else {
		startDaemon()
	}
}
