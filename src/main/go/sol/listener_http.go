package main

import (
	"encoding/xml"
	"io"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type route struct {
	pattern *regexp.Regexp
	handler func(http.ResponseWriter, *http.Request)
}

type RegexpHandler struct {
	routes []*route
}

type RestResultHost struct {
	XMLName    xml.Name `xml:"host"`
	Ip         string   `xml:"ip,attr"`
	MacAddress string   `xml:"mac,attr"`
}

type RestResultHosts struct {
	XMLName xml.Name `xml:"hosts"`
	Hosts   []RestResultHost
}

type RestResultCommands struct {
	XMLName  xml.Name `xml:"commands"`
	Commands []RestResultCommandConfiguration
}

type RestResultCommandConfiguration struct {
	XMLName   xml.Name `xml:"command"`
	Operation string   `xml:"operation,attr"`
	Command   string   `xml:"command,attr"`
	IsDefault bool     `xml:"default,attr"`
	Type      string   `xml:"type,attr"`
}

type RestResultListeners struct {
	XMLName   xml.Name `xml:"listeners"`
	Listeners []RestResultListenerConfiguration
}

type RestResultListenerConfiguration struct {
	XMLName xml.Name `xml:"listener"`
	Type    string   `xml:"type,attr"`
	Port    int      `xml:"port,attr"`
	Active  bool     `xml:"active,attr"`
}

type RestResult struct {
	XMLName              xml.Name `xml:"result"`
	Application          string   `xml:"application"`
	Version              string   `xml:"version"`
	CompilationTimestamp string   `xml:"compilation"`
	Hosts                RestResultHosts
	Listeners            RestResultListeners
	Commands             RestResultCommands
}

func (h *RegexpHandler) Handler(re string, handler func(http.ResponseWriter, *http.Request)) {
	r := &route{regexp.MustCompile(re), handler}
	h.routes = append(h.routes, r)
}

// func (h *RegexpHandler) HandleFunc(pattern *regexp.Regexp, handler func(http.ResponseWriter, *http.Request)) {
//     h.routes = append(h.routes, &route{pattern, http.HandlerFunc(handler)})
//}

func (h *RegexpHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	for _, route := range h.routes {
		if route.pattern.MatchString(r.URL.Path) {
			route.handler(rw, r)
			return
		}
	}
	http.NotFound(rw, r)
}

func restGenericOperation(w http.ResponseWriter, r *http.Request) {
	items := strings.Split(r.URL.Path, "/")
	operation := items[1]

	for _, Command := range configuration.Commands {
		if Command.Operation == operation {
			io.WriteString(w, "Executing operation ["+operation+"]")
			ExecuteCommand(Command)
			break
		}
	}
}

func restIndex(w http.ResponseWriter, r *http.Request) {
	result := &RestResult{}
	result.Application = APPLICATION_NAME
	result.Version = VERSION
	if BUILD_DATE != "" {
		result.CompilationTimestamp = BUILD_DATE
	}
	result.Hosts = RestResultHosts{}
	result.Listeners = RestResultListeners{}
	result.Commands = RestResultCommands{}

	interfaces := LocalNetworkMap()
	ips := make([]string, 0, len(interfaces))
	for key := range interfaces {
		ips = append(ips, key)
	}
	sort.Strings(ips)
	for _, ip := range ips {
		result.Hosts.Hosts = append(result.Hosts.Hosts, RestResultHost{Ip: ip, MacAddress: interfaces[ip]})
	}
	for _, listenerConfiguration := range configuration.listenersConfiguration {
		result.Listeners.Listeners = append(result.Listeners.Listeners, RestResultListenerConfiguration{Type: listenerConfiguration.nature, Port: listenerConfiguration.port, Active: listenerConfiguration.active})
	}

	for _, commandConfiguration := range configuration.Commands {
		result.Commands.Commands = append(result.Commands.Commands, RestResultCommandConfiguration{Type: commandConfiguration.CommandType, Operation: commandConfiguration.Operation, Command: commandConfiguration.Command, IsDefault: commandConfiguration.IsDefault})
	}

	x, err := xml.MarshalIndent(result, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/xml")
	w.Write(x)
}

func restStopSleepOnLan(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Now quitting sleep on lan")
	panic("Quit")
}

func restWol(w http.ResponseWriter, r *http.Request) {
	mac := strings.Replace(r.URL.Path, "/wol/", "", 1)
	Info.Println("Now sending wol magic packet to MAC address [" + mac + "]")
	magicPacket, err := EncodeMagicPacket(mac)
	if err != nil {
		Error.Println(err)
	}

	magicPacket.Wake(configuration.BroadcastIP)
}

func ListenHTTP(port int, commands []CommandConfiguration) {
	localIp := "0.0.0.0"

	routes := make(map[string]func(http.ResponseWriter, *http.Request))
	// routes["/sleep"] = restSleep
	routes["/wol/[A-z]+"] = restWol
	for _, command := range commands {
		routes["/"+command.Operation] = restGenericOperation
	}
	routes["/"] = restIndex

	// routes["/quit"] = restStopSleepOnLan

	externalIp, _ := ExternalIP()
	baseExternalUrl := "http://" + externalIp + ":" + strconv.Itoa(port)

	reHandler := new(RegexpHandler)
	Info.Println("Now listening HTTP on port [" + strconv.Itoa(port) + "], urls will be : ")
	for key, value := range routes {
		reHandler.Handler(key, value)
		Info.Println(" - " + baseExternalUrl + key)
	}
	http.ListenAndServe(localIp+":"+strconv.Itoa(port), reHandler)
}
