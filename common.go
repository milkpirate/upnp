package upnp

import (
	// "log"
	"errors"
	"net"
	"strings"
)

// GetLocalIntenetIp returns all local addresses,
// checks internet connectivity
func GetLocalIntenetIp() string {
	/*
	  Get all local addresses
	  Check for internet connectivity
	*/

	conn, err := net.Dial("udp4", "google.com:80")
	if err != nil {
		panic(errors.New("Can not connect to network"))
	}
	defer conn.Close()
	return strings.Split(conn.LocalAddr().String(), ":")[0]
}

// This returns the list of local ip addresses which other hosts can connect
// to (NOTE: Loopback ip is ignored).
func GetLocalIPs() ([]*net.IP, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	ips := make([]*net.IP, 0)
	for _, addr := range addrs {
		ipnet, ok := addr.(*net.IPNet)
		if !ok {
			continue
		}

		if ipnet.IP.IsLoopback() {
			continue
		}

		ips = append(ips, &ipnet.IP)
	}

	return ips, nil
}
