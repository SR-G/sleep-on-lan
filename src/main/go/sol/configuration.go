package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type Configuration struct {
	Listeners   []string // what is read from the sol.json configuration file
	LogLevel    string
	BroadcastIP string
	Commands    []CommandConfiguration // the various defined commands. Will be enhanded with default operation if empty from configuration

	listenersConfiguration []ListenerConfiguration // converted once parsed from Listeners
}

type CommandConfiguration struct {
	Operation   string `json:"Operation"`
	Command     string `json:"Command"`
	IsDefault   bool   `json:"Default"`
	CommandType string `json:"Type"`
}

type ListenerConfiguration struct {
	active bool
	port   int
	nature string
}

func (conf *Configuration) InitDefaultConfiguration() {
	conf.Listeners = []string{"UDP:9", "HTTP:8009"}
	conf.LogLevel = "INFO"
	conf.BroadcastIP = "192.168.255.255"
	// default commands are registered on Parse() method, depending on the current operating system
}

func (conf *Configuration) Load(configurationFileName string) {
	if _, err := os.Stat(configurationFileName); err == nil {
		Info.Println("Configuration file found under [" + configurationFileName + "] now reading content")
		file, _ := os.Open(configurationFileName)
		decoder := json.NewDecoder(file)
		err := decoder.Decode(&conf)
		if err != nil {
			Error.Println("error while loading configuration :", err)
		}
	} else {
		Info.Println("No external configuration file found under [" + configurationFileName + "] will use default values")
	}
}

func (conf *Configuration) Parse() {
	// Gestion logs
	switch conf.LogLevel {
	case "NONE", "OFF":
		InitLoggers(ioutil.Discard, ioutil.Discard, ioutil.Discard, ioutil.Discard)
	case "DEBUG":
		InitLoggers(os.Stdout, os.Stdout, os.Stdout, os.Stderr)
	case "INFO":
		InitLoggers(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)
	case "WARN", "WARNING":
		InitLoggers(ioutil.Discard, ioutil.Discard, os.Stdout, os.Stderr)
	case "ERROR":
		InitLoggers(ioutil.Discard, ioutil.Discard, ioutil.Discard, os.Stderr)
	default:
		panic("unrecognized log level[" + conf.LogLevel + "], allowed are NONE or OFF, DEBUG, INFO, WARN or WARNING, ERROR")
	}

	// Convert activated ports
	for _, s := range conf.Listeners {
		var splitted = strings.Split(s, ":")
		var key = splitted[0]
		var listenerConfiguration = new(ListenerConfiguration)
		listenerConfiguration.active = true
		if len(splitted) == 2 {
			listenerConfiguration.port, _ = strconv.Atoi(splitted[1])
		}
		if strings.EqualFold(key, "UDP") {
			listenerConfiguration.nature = "UDP"
			conf.listenersConfiguration = append(conf.listenersConfiguration, *listenerConfiguration)
		} else if strings.EqualFold(key, "HTTP") {
			listenerConfiguration.nature = "HTTP"
			conf.listenersConfiguration = append(conf.listenersConfiguration, *listenerConfiguration)
		} else {
			Error.Println("Unknown listener type [" + key + "], valid values are : UDP, HTTP")
		}
	}
	Trace.Println("Configuration loaded", conf)

	// If only one command, then force default, and if no commands are found, inject default ones
	var nbCommands = len(conf.Commands)
	if nbCommands == 0 {
		RegisterDefaultCommand()
	} else if nbCommands == 1 {
		Info.Println("Only one command found in configuration, forcing default if needed")
		conf.Commands[0].IsDefault = true
	}
}
