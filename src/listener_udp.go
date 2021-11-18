package main

import (
	"net"
	"strconv"
	"strings"
	"time"
)

type MagicPacket []byte

var isActionInProgress bool = false

func ListenUDP(port int) {
	logger.Infof("Now listening UDP packets on port [" + colorer.Green(strconv.Itoa(port)) + "]")
	var buf [1024]byte
	addr, err := net.ResolveUDPAddr("udp", ":"+strconv.Itoa(port))
	if err != nil {
		logger.Errorf("Error while resolving local address :", err.Error())
	}
	sock, err := net.ListenUDP("udp", addr)
	if err != nil {
		if configuration.ExitIfAnyPortIsAlreadyUsed {
			logger.Errorf("Error while starting listening (will exit, per configuration) :", err.Error())
			exit <- true
			return
		} else {
			logger.Errorf("Error while starting listening (program will however continue) :", err.Error())
		}
	} else {
		for {
			rlen, remote, err := sock.ReadFromUDP(buf[:])
			if err != nil {
				logger.Errorf("Error while reading :", err.Error())
			}
			extractedMacAddress, _ := extractMacAddress(rlen, buf)
			logger.Infof("Received a MAC address from IP [" + remote.String() + "], extracted mac [" + extractedMacAddress.String() + "]")
			if matchAddress(extractedMacAddress) {
				logger.Infof("(reversed) received MAC address match a local address")
				if configuration.AvoidDualUDPSending.AvoidDualUDPSendingActive {
					// Specific behavior : let's try to avoid dual UDP sending
					if !isActionInProgress {
						isActionInProgress = true
						logger.Infof("Extra small delay before going to sleep (to avoid dual UDP sending), during [" + configuration.AvoidDualUDPSending.AvoidDualUDPSendingDelay + "]")
						go doActionWithDelay()
					} else {
						logger.Infof("Another command is already awaiting, rejecting this one due to dual UDP sending avoidance being activated")
					}
				} else {
					// Regular behavior, let's just execute command
					doAction()
				}
			}
		}
	}
}

func doActionWithDelay() {
	delay, _ := time.ParseDuration(configuration.AvoidDualUDPSending.AvoidDualUDPSendingDelay)
	time.Sleep(delay)
	doAction()
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
		logger.Errorf("Received buffer too small, size [" + strconv.Itoa(rlen) + "]")
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
	for idx, _ := range configuration.Commands {
		Command := configuration.Commands[idx]
		if Command.IsDefault {
			isActionInProgress = false
			ExecuteCommand(Command)
			break
		}
	}
}
