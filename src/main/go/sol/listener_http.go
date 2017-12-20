package main

import (
	"encoding/xml"
	"net/http"
	"sort"
	"strconv"
	// "encoding/json"
	// "io/ioutil"
	"os"
	"strings"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	// "github.com/labstack/echo/engine/standard"
)

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

type RestOperationResult struct {
	XMLName              xml.Name `xml:"result"`
	Operation            string   `xml:"operation"`
	Result				 bool     `xml:"successful"`
}

func dumpRoute(route string) {
	Info.Println("Registering route [/" + route + "]")
}

func ListenHTTP(port int, commands []CommandConfiguration, auth AuthConfiguration) {
	// externalIp, _ := ExternalIP()
	// baseExternalUrl := "http://" + externalIp + ":" + strconv.Itoa(port)
	// Info.Println("Now listening HTTP on port [" + strconv.Itoa(port) + "], urls will be : ")
	/*
	for key, value := range routes {
		Info.Println(" - " + baseExternalUrl + key)
	}
	*/

	e := echo.New()
	e.HideBanner = true
	
	// e.File("/", "public/index.html")
	// e.Static("/", "public")
	// e.Use(middleware.Gzip())

	if auth.isEmpty() {
		Info.Println("HTTP starting on port [" + strconv.Itoa(port) + "], without auth")
	} else {
		Info.Println("HTTP starting on port [" + strconv.Itoa(port) + "], with auth activated : login [" + auth.Login + "], password [" + strings.Repeat("*", len(auth.Password)) + "]")
		e.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
			if username == auth.Login && password == auth.Password {
				return true, nil
			}
			return false, nil
		}))
	}

	dumpRoute("")
	e.GET("/", func(c echo.Context) error {
		result := &RestResult{}
		result.Application = Version.VersionLabel
		result.Version = Version.Version()
		if Build != "" {
			result.CompilationTimestamp = Build
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
	
		return c.XMLPretty(http.StatusOK, result, "  ")
	})

	// N.B.: sleep operation is now registred through commands below
	for _, command := range commands {
		dumpRoute(command.Operation)
		e.GET("/" + command.Operation, func(c echo.Context) error {
			
			items := strings.Split(c.Request().URL.Path, "/")
			operation := items[1]

			result := &RestOperationResult{
				Operation:  operation,
				Result: true,
			}
			for idx, _ := range configuration.Commands {
				availableCommand := configuration.Commands[idx]
				if availableCommand.Operation == operation {
					Info.Println("Executing [" + operation + "]")
					ExecuteCommand(availableCommand)
					break
				}
		}			
		return c.XMLPretty(http.StatusOK, result, "  ")
		})
	}

	dumpRoute("quit")
	e.GET("/quit", func(c echo.Context) error {
		result := &RestOperationResult{
			Operation:  "quit",
			Result: true,
		}
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextXMLCharsetUTF8)
		// c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationXMLCharsetUTF8) // echo.MIMETextXMLCharsetUTF8)
		c.Response().WriteHeader(http.StatusOK)		
		b, _ := xml.Marshal(result)
		c.Response().Write(b)
		c.Response().Flush()
		defer os.Exit(1)
		return nil
		// return c.XMLPretty(http.StatusOK, result, "  ")
	})

	dumpRoute("wol/:mac")
	e.GET("/wol/:mac", func(c echo.Context) error {
		result := &RestOperationResult{
			Operation:  "wol",
			Result: true,
		}

        mac := c.Param("mac")
		Info.Println("Now sending wol magic packet to MAC address [" + mac + "]")
		magicPacket, err := EncodeMagicPacket(mac)
		if err != nil {
			Error.Println(err)
		} else {
			magicPacket.Wake(configuration.BroadcastIP)
		}
		return c.XMLPretty(http.StatusOK, result, "  ")
	})
	
	// localIp := "0.0.0.0"
	Info.Println(e.Start(":" + strconv.Itoa(port)))
}