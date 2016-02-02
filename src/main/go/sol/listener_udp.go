package main

import (
	"net"
	"strconv"
	"strings"
)

type MagicPacket []byte

func ListenUDP(port int) {
	Info.Println("Now listening UDP packets on port [" + strconv.Itoa(port) + "]")
	var buf [1024]byte
	addr, err := net.ResolveUDPAddr("udp", ":"+strconv.Itoa(port))
	if err != nil {
		Error.Println("Error while resolving local address :", err.Error())
	}
	sock, err := net.ListenUDP("udp", addr)
	if err != nil {
		Error.Println("Error while starting listening :", err.Error())
	}
	for {
		rlen, remote, err := sock.ReadFromUDP(buf[:])
		if err != nil {
			Error.Println("Error while reading :", err.Error())
		}
		extractedMacAddress, _ := extractMacAddress(rlen, buf)
		Info.Println("Received a MAC address from IP [" + remote.String() + "], extracted mac [" + extractedMacAddress.String() + "]")
		if matchAddress(extractedMacAddress) {
			Info.Println("(reversed) received MAC address match a local address")
			doAction()
		}
	}
}

func matchAddress(receivedAddress net.HardwareAddr) bool {
	receivedAddressAsString := receivedAddress.String()
	for _, value := range LocalNetworkMap() {
		if strings.HasPrefix(value, receivedAddressAsString) {
			return true
		}
		/*if bytes.Equal(receivedAddress, inter.HardwareAddr) {
			return true
		}*/
	}

	return false
}

func extractMacAddress(rlen int, buf [1024]byte) (net.HardwareAddr, error) {
	var r = ""
	// TODO check whole magic packet structure (FF FF FF FF FF FF <MAC>*6)
	if rlen >= 12 {
		var sep = ""
		for i := 6; i < 12; i++ {
			val := int64(buf[i])                 // decimal value
			s := strconv.FormatInt(val, 16)      // convert to hexa (base 16)
			r = leftPad2Len(s, "0", 2) + sep + r // pad on two characters because some wake on lan tools are actually sending ":01:" as ":1:"
			sep = ":"
		}
	} else {
		Error.Println("Received buffer too small, size [" + strconv.Itoa(rlen) + "]")
	}
	return net.ParseMAC(r)
}

func leftPad2Len(s string, padStr string, overallLen int) string {
	var padCountInt int
	padCountInt = 1 + ((overallLen - len(padStr)) / len(padStr))
	var retStr = strings.Repeat(padStr, padCountInt) + s
	return retStr[(len(retStr) - overallLen):]
}

func doAction() {
	for _, Command := range configuration.Commands {
		if Command.IsDefault {
			ExecuteCommand(Command)
			break
		}
	}
}
