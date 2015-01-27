package main

import (
	"encoding/xml"
	"io"
	"net/http"
	"regexp"
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
	XMLName     xml.Name `xml:"result"`
	Application string   `xml:"application"`
	Hosts       RestResultHosts
	Listeners   RestResultListeners
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

func restSleep(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Going to sleep now")
	Sleep()
}

func restIndex(w http.ResponseWriter, r *http.Request) {
	result := &RestResult{}
	result.Application = "sleep-on-lan"
	result.Hosts = RestResultHosts{}
	result.Listeners = RestResultListeners{}
	for key, value := range LocalNetworkMap() {
		result.Hosts.Hosts = append(result.Hosts.Hosts, RestResultHost{Ip: key, MacAddress: value})
	}
	for _, listenerConfiguration := range configuration.listenersConfiguration {
		result.Listeners.Listeners = append(result.Listeners.Listeners, RestResultListenerConfiguration{Type: listenerConfiguration.nature, Port: listenerConfiguration.port, Active: listenerConfiguration.active})
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

func ListenHTTP(port int) {
	localIp := "0.0.0.0"

	routes := make(map[string]func(http.ResponseWriter, *http.Request))
	routes["/sleep"] = restSleep
	routes["/wol/[A-z]+"] = restWol
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
