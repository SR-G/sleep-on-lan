package main

import (
	"encoding/json"
	"fmt"
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

func ExitDaemon(s string) {
	if s == "" {
		logger.Errorf("Stopping daemon ...")
	} else {
		logger.Errorf("Stopping daemon, " + s)
	}
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
	// Prepare new default configuration and associated bytes
	c := Configuration{}
	c.InitDefaultConfiguration()
	c.Parse()

	b, err := json.MarshalIndent(c, "", "    ")
	if err != nil {
		logger.Errorf("Can't create JSON content for default content", err)
	} else {
		// Check which configuration file to use (either from --config, either sol.json file alongside the binary, either default values)
		if configurationFileNameFromCommandLine == "" {
			// Write bytes on console (default behavior when -c is not used)
			fmt.Println("Example of possible default configuration (to be stored alongside sol binary in a filename like [" + colorer.Green("sol.json") + "])")
			fmt.Print(string(b))
		} else {
			// Write bytes on disk (if possible)
			if _, err := os.Stat(configurationFileNameFromCommandLine); err == nil {
				logger.Errorf("Can't generate default JSON configuration info [" + colorer.Red(configurationFileNameFromCommandLine) + "] : file already exist")
			} else {
				logger.Infof("Writing default JSON configuration into [" + colorer.Green(configurationFileNameFromCommandLine) + "]")

				file, err := os.OpenFile(configurationFileNameFromCommandLine, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
				if err != nil {
					logger.Errorf("Can't write file ["+colorer.Red(configurationFileNameFromCommandLine)+"]", err)
				}
				defer file.Close()

				// Create blank configuration (with default values) and write bytes to file
				bytesWritten, err := file.Write(b)
				if err != nil {
					logger.Errorf("Can't write JSON content into ["+colorer.Red(configurationFileNameFromCommandLine)+"]", err)
				} else {
					logger.Infof("Bytes written : " + strconv.Itoa(bytesWritten))
				}
			}
		}
	}
}

func determineConfigurationFileName() string {
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
	return fullConfigurationFileName
}

func startDaemon() {
	// Check which configuration file to use (either from --config, either sol.json file alongside the binary, either default values)
	fullConfigurationFileName := determineConfigurationFileName()

	// Loads configuration, and continue only if all configuration entries are OK
	configurationOk := true
	configurationError := ""
	configuration.InitDefaultConfiguration()
	if err := configuration.Load(fullConfigurationFileName); err != nil {
		configurationOk = false
		configurationError = err.Error()
	}
	if err := configuration.RefineLogger(); err != nil {
		configurationOk = false
		configurationError = err.Error()
	}
	if err := configuration.Parse(); err != nil {
		configurationOk = false
		configurationError = err.Error()
	}

	if !configurationOk {
		logger.Errorf("Stopping daemon due to configuration errors : " + colorer.Red(configurationError))
	} else {
		logger.Infof("Application [" + colorer.Green(Version.ApplicationName) + "], version [" + colorer.Green(Version.GetVersion()) + "], compilation timestamp [" + colorer.Green(Version.CompilationTimestamp) + "], git commit [" + colorer.Green(Version.Commit) + "]")

		// Display found IP/MAC
		logger.Infof("Now starting sleep-on-lan, hardware IP/mac addresses are : ")
		for key, value := range LocalNetworkMap() {
			logger.Infof(" - local IP adress [" + colorer.Green(key) + "], mac [" + colorer.Green(value) + "], reversed mac [" + colorer.Green(ReverseMacAddress(value)) + "]")
		}

		// Display commands found in configuration
		logger.Infof("Available commands are : ")
		for _, command := range configuration.Commands {
			if command.Command == "" {
				logger.Infof(" - operation [" + color.Green(command.Operation) + "], default [" + color.Green(strconv.FormatBool(command.IsDefault)) + "], type [" + color.Green(command.CommandType) + "]")
			} else {
				logger.Infof(" - operation [" + color.Green(command.Operation) + "], command [" + color.Green(command.Command) + "], default [" + color.Green(strconv.FormatBool(command.IsDefault)) + "], type [" + color.Green(command.CommandType) + "]")
			}
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
}

func main() {
	// Handling subcommands, if any
	if subCommandGenerateConfiguration.Used {
		executeCommandGenerateConfiguration()
	} else {
		startDaemon()
	}
}
