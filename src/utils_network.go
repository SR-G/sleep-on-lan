package main

import (
	"encoding/hex"
	"errors"
	"net"
	"strings"
)

func ExternalIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String(), nil
		}
	}
	return "", errors.New("are you connected to the network?")
}

// Use the net library to return all Interfaces
// and capture any errors.
func GetInterfaces() []net.Interface {
	interfaces, err := net.Interfaces()
	if err != nil {
		logger.Errorf("Unable to get interfaces : " + err.Error())
	}
	return interfaces
}

// Use the net library to get Interface by Index number
// and capture any errors.
func GetInterfaceByIndex(index int) (*net.Interface, error) {
	iface, err := net.InterfaceByIndex(index)
	if err != nil {
		return nil, errors.New("Unable to get interface by index : " + err.Error())
	}
	return iface, nil
}

// Use the net library to get Interface by Name
// and capture any errors
func GetInterfaceByName(name string) *net.Interface {
	iface, err := net.InterfaceByName(name)
	if err != nil {
		panic("Unable to get interface by name")
	}
	return iface
}

// Using the net library we will loop over all interfaces
// looking for the interface with matching mac_address.
func GetInterfaceByHardwareAddress(mac string) (net.Interface, error) {
	interfaces := GetInterfaces()
	for _, iface := range interfaces {
		if strings.EqualFold(mac, iface.HardwareAddr.String()) {
			return iface, nil
		}
	}
	panic("Unable to find interface with Hardware Address.")
}

// Use a MAC address to form a magic packet
// macAddr form 12:34:56:78:9a:bc
func EncodeMagicPacket(macAddr string) (MagicPacket, error) {
	if len(macAddr) != (6*2 + 5) {
		return nil, errors.New("Invalid MAC Address [" + macAddr + "]")
	}

	macBytes, err := hex.DecodeString(strings.Join(strings.Split(macAddr, ":"), ""))
	if err != nil {
		return nil, err
	}

	b := []uint8{255, 255, 255, 255, 255, 255}
	for i := 0; i < 16; i++ {
		b = append(b, macBytes...)
	}

	return MagicPacket(b), nil
}

func LocalNetworkMap() map[string]string {
	result := make(map[string]string)
	for _, inter := range GetInterfaces() {
		addresses, _ := inter.Addrs()
		for _, addr := range addresses {
			result[addr.String()] = inter.HardwareAddr.String()
		}
		// logger.Infof(inter.Name , " : ", inter.HardwareAddr, ", ")
		// logger.Infof(inter.Addrs())
	}
	return result
}

// Send a Magic Packet to an broadcast class IP address via UDP
func (p MagicPacket) Wake(bcastAddr string) error {
	a, err := net.ResolveUDPAddr("udp", bcastAddr+":9")
	if err != nil {
		return err
	}

	c, err := net.DialUDP("udp", nil, a)
	if err != nil {
		return err
	}

	written, err := c.Write(p)
	c.Close()

	// Packet must be 102 bytes in length
	if written != 102 {
		return err
	}

	return nil
}
