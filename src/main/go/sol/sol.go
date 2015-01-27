package main

import (
	"os"
	"path/filepath"
	"strings"
)

var configuration = Configuration{}
var configurationFileName = "sol.json"

func main() {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	PreInitLoggers()
	configuration.InitDefaultConfiguration()
	configuration.Load(dir + string(os.PathSeparator) + configurationFileName)
	configuration.Parse()

	Info.Println("Now starting sleep-on-lan, hardware IP/mac addresses are : ")
	for key, value := range LocalNetworkMap() {
		Info.Println(" - local IP adress [" + key + "], mac [" + value + "]")
	}

	for _, listenerConfiguration := range configuration.listenersConfiguration {
		if listenerConfiguration.active {
			if strings.EqualFold(listenerConfiguration.nature, "UDP") {
				go ListenUDP(listenerConfiguration.port)
			} else if strings.EqualFold(listenerConfiguration.nature, "HTTP") {
				go ListenHTTP(listenerConfiguration.port)
			}
		}
	}

	// Info.Println("sleep-on-lan up and running")
	select {} // block forever
}
