package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type Configuration struct {
	Listeners              []string
	SleepCommand           string
	LogLevel               string
	BroadcastIP            string
	listenersConfiguration []ListenerConfiguration
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
	Info.Println("Configuration loaded", conf)

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
}
