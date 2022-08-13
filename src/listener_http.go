package main

import (
	"encoding/xml"
	"net/http"

	// "os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"github.com/go-ping/ping" // for ping
	// "github.com/mdlayher/arp" // for mac > ip conversion
)

type RestResultHost struct {
	XMLName            xml.Name `xml:"host" json:"-"`
	Ip                 string   `xml:"ip,attr"`
	MacAddress         string   `xml:"mac,attr"`
	ReversedMacAddress string   `xml:"reversed-mac,attr" json:"-"`
}

type RestResultHosts struct {
	XMLName xml.Name `xml:"hosts" json:"-"`
	Hosts   []RestResultHost
}

type RestResultCommands struct {
	XMLName  xml.Name `xml:"commands" json:"-"`
	Commands []RestResultCommandConfiguration
}

type RestResultCommandConfiguration struct {
	XMLName   xml.Name `xml:"command" json:"-"`
	Operation string   `xml:"operation,attr"`
	Command   string   `xml:"command,attr"`
	IsDefault bool     `xml:"default,attr"`
	Type      string   `xml:"type,attr"`
}

type RestResultListeners struct {
	XMLName   xml.Name `xml:"listeners" json:"-"`
	Listeners []RestResultListenerConfiguration
}

type RestResultListenerConfiguration struct {
	XMLName xml.Name `xml:"listener" json:"-"`
	Type    string   `xml:"type,attr"`
	Port    int      `xml:"port,attr"`
	Active  bool     `xml:"active,attr"`
}

type RestResult struct {
	XMLName              xml.Name `xml:"result" json:"-"`
	Application          string   `xml:"application"`
	Version              string   `xml:"version,omitifempty"`
	CompilationTimestamp string   `xml:"compilation-timestamp,omitifempty"`
	Commit               string   `xml:"commit,omitifempty"`
	Hosts                RestResultHosts
	Listeners            RestResultListeners
	Commands             RestResultCommands
}

type RestOperationResult struct {
	XMLName   xml.Name `xml:"result" json:"-"`
	Operation string   `xml:"operation"`
	Result    bool     `xml:"successful"`
}

const (
	HOST_STATE_ONLINE  = "online"
	HOST_STATE_OFFLINE = "offline"
	HOST_STATE_UNKNOWN = "unknown"
)

type RestStateResult struct {
	XMLName xml.Name `xml:"result" json:"-"`
	State   string   `xml:"state"`
	Host    string   `xml:"host"`
}

func dumpRoute(route string) {
	logger.Infof("Registering route [" + colorer.Green("/"+route) + "]")
}

// func retrieveIpFromMac(mac strinc) string {
// requires defined interface ...
// }

func renderResult(c echo.Context, status int, result interface{}) error {
	// c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationXMLCharsetUTF8) // echo.MIMETextXMLCharsetUTF8)

	// Return status cope
	c.Response().WriteHeader(status)

	// Proper formatting per what is expected in the query
	format := c.QueryParam("format")
	if strings.EqualFold(configuration.HTTPOutput, "JSON") || strings.EqualFold(format, "JSON") {
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		return c.JSONPretty(status, result, "  ")
	} else {
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextXMLCharsetUTF8)
		return c.XMLPretty(status, result, "  ")
	}
	// c.Response().Write(b)
	// c.Response().Flush()
}

func pingIp(ip string) *RestStateResult {
	logger.Infof("Checking state of remote host with IP [" + ip + "]")
	result := &RestStateResult{
		Host:  ip,
		State: HOST_STATE_ONLINE,
	}
	pinger, err := ping.NewPinger(ip)
	if err != nil {
		logger.Infof("Can't retrieve PING results (rights problems when executing sol, maybe ?)", err)
		result.State = HOST_STATE_UNKNOWN
	}
	pinger.Count = 3
	// pinger.Interval = // default is 1s, which is fine
	pinger.Timeout = time.Second * 5
	pinger.SetPrivileged(true)
	pinger.Run()                                                                                                                                                                                                                                              // blocks until finished
	stats := pinger.Statistics()                                                                                                                                                                                                                              // get send/receive/rtt stats
	logger.Infof("Ping results for [" + stats.Addr + "], [" + strconv.Itoa(stats.PacketsSent) + "] packets transmitted, [" + strconv.Itoa(stats.PacketsRecv) + "] packets received, [" + strconv.FormatFloat(stats.PacketLoss, 'f', 2, 64) + "] packet loss") // , round-trip min/avg/max/stddev = " + stats.MinRtt + "/" + stats.AvgRtt + "/" + stats.MaxRtt + "/" + stats.StdDevRtt + "")
	if stats.PacketsRecv == 0 {
		result.State = HOST_STATE_OFFLINE
	}
	return result
}

func ExecuteActionWithDelay(availableCommand CommandConfiguration) {
	delay, _ := time.ParseDuration(configuration.DelayBeforeCommands.Delay)
	time.Sleep(delay)
	ExecuteCommand(availableCommand)
}

func ListenHTTP(port int) {
	// externalIp, _ := ExternalIP()
	// baseExternalUrl := "http://" + externalIp + ":" + strconv.Itoa(port)
	// logger.Infof("Now listening HTTP on port [" + strconv.Itoa(port) + "], urls will be : ")
	/*
		for key, value := range routes {
						logger.Infof(" - " + baseExternalUrl + key)
								}
	*/

	e := echo.New()
	e.HideBanner = true

	// e.File("/", "public/index.html")
	// e.Static("/", "public")
	// e.Use(middleware.Gzip())

	if configuration.Auth.isEmpty() {
		logger.Infof("HTTP starting on port [" + colorer.Green(strconv.Itoa(port)) + "], without auth")
	} else {
		logger.Infof("HTTP starting on port [" + colorer.Green(strconv.Itoa(port)) + "], with auth activated : login [" + configuration.Auth.Login + "], password [" + strings.Repeat("*", len(configuration.Auth.Password)) + "]")
		e.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
			if username == configuration.Auth.Login && password == configuration.Auth.Password {
				return true, nil
			}
			return false, nil
		}))
	}

	dumpRoute("")
	e.GET("/", func(c echo.Context) error {
		result := &RestResult{}
		result.Application = Version.ApplicationName
		result.Version = Version.GetVersion()
		result.Commit = Version.Commit
		if Version.CompilationTimestamp != "" {
			result.CompilationTimestamp = Version.CompilationTimestamp
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
			result.Hosts.Hosts = append(result.Hosts.Hosts, RestResultHost{Ip: ip, MacAddress: interfaces[ip], ReversedMacAddress: ReverseMacAddress(interfaces[ip])})
		}
		for _, listenerConfiguration := range configuration.listenersConfiguration {
			result.Listeners.Listeners = append(result.Listeners.Listeners, RestResultListenerConfiguration{Type: listenerConfiguration.nature, Port: listenerConfiguration.port, Active: listenerConfiguration.active})
		}

		for _, commandConfiguration := range configuration.Commands {
			result.Commands.Commands = append(result.Commands.Commands, RestResultCommandConfiguration{Type: commandConfiguration.CommandType, Operation: commandConfiguration.Operation, Command: commandConfiguration.Command, IsDefault: commandConfiguration.IsDefault})
		}

		return renderResult(c, http.StatusOK, result)
	})

	// N.B.: sleep operation is now registred through commands below
	for _, command := range configuration.Commands {
		dumpRoute(command.Operation)
		e.GET("/"+command.Operation, func(c echo.Context) error {

			items := strings.Split(c.Request().URL.Path, "/")
			operation := items[1]

			result := &RestOperationResult{
				Operation: operation,
				Result:    true,
			}
			for idx := range configuration.Commands {
				availableCommand := configuration.Commands[idx]
				if availableCommand.Operation == operation {

					if configuration.DelayBeforeCommands.Active {
						logger.Infof("Executing [" + colorer.Green(operation) + "] with preliminary delay of [" + colorer.Green(configuration.DelayBeforeCommands.Delay) + "]")
						// we still need "go" in order to not be blocking here...
						go ExecuteActionWithDelay(availableCommand)
					} else {
						logger.Infof("Executing [" + colorer.Green(operation) + "]")
						go ExecuteCommand(availableCommand)
					}

					break
				}
			}
			return renderResult(c, http.StatusOK, result)
		})
	}

	dumpRoute("quit")
	e.GET("/quit", func(c echo.Context) error {
		result := &RestOperationResult{
			Operation: "quit",
			Result:    true,
		}
		defer ExitDaemon("quit command triggered")
		return renderResult(c, http.StatusOK, result)

	})

	dumpRoute("state/local/online")
	e.GET("/state/local/online", func(c echo.Context) error {
		return c.String(http.StatusOK, "true")
	})

	dumpRoute("state/local")
	e.GET("/state/local", func(c echo.Context) error {
		result := &RestStateResult{
			Host:  "localhost",
			State: HOST_STATE_ONLINE,
		}
		return renderResult(c, http.StatusOK, result)
	})

	dumpRoute("state/ip/:ip")
	e.GET("/state/ip/:ip", func(c echo.Context) error {
		ip := c.Param("ip")
		result := pingIp(ip)
		return renderResult(c, http.StatusOK, result)
	})

	/*
		dumpRoute("state/mac/:mac")
		e.GET("/state/ip/:ip", func(c echo.Context) error {
		mac := c.Param("mac")
		ip := retrieveIpFromMac(mac)
		result := pingIp(ip)
		return c.XMLPretty(http.StatusOK, result, "  ")
	*/

	dumpRoute("wol/:mac")
	e.GET("/wol/:mac", func(c echo.Context) error {
		result := &RestOperationResult{
			Operation: "wol",
			Result:    true,
		}

		mac := c.Param("mac")
		logger.Infof("Now sending wol magic packet to MAC address [" + mac + "]")
		magicPacket, err := EncodeMagicPacket(mac)
		if err != nil {
			logger.Errorf("Error while sending magic packet to MAC address ["+mac+"]", err)
		} else {
			magicPacket.Wake(configuration.BroadcastIP)
		}
		return renderResult(c, http.StatusOK, result)
	})

	// localIp := "0.0.0.0"
	if err := e.Start(":" + strconv.Itoa(port)); err != http.ErrServerClosed {
		if configuration.ExitIfAnyPortIsAlreadyUsed {
			logger.Errorf("Unable to start HTTP listener on port [" + strconv.Itoa(port) + "] (program will be stopped, per configuration) : " + err.Error())
			defer ExitDaemon("port already taken")
		} else {
			logger.Errorf("Unable to start HTTP listener on port [" + strconv.Itoa(port) + "] (program will be continue) : " + err.Error())
		}
	}
}
