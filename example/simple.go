package main

import (
	// "bufio"
	"fmt"

	"github.com/scottjg/upnp"
	// "os"
)

func main() {
	SearchGateway()
	ExternalIPAddr()
	AddPortMapping()
}

// Search gateway device
func SearchGateway() {
	upnpMan := new(upnp.Upnp)
	err := upnpMan.SearchGateway()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Local IP Address：", upnpMan.LocalHost)
		fmt.Println("UPNP Device IP Address:", upnpMan.Gateway.Host)
	}
}

// Obtain public ip address
func ExternalIPAddr() {
	upnpMan := new(upnp.Upnp)
	err := upnpMan.ExternalIPAddr()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("External Network IP Address:", upnpMan.GatewayOutsideIP)
	}
}

// Add a port mapping
func AddPortMapping() {
	mapping := new(upnp.Upnp)
	if err := mapping.AddPortMapping(55789, 55789, 0, "TCP"); err == nil {
		fmt.Println("Port mapping succeeded.")
		mapping.Reclaim()
	} else {
		fmt.Println("Port mapping failed.")
	}

}
