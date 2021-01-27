package main

import (
	// "bufio"
	"fmt"

	"github.com/fdurand/upnp"
	// "os"
)

func main() {
	SearchGateway()
	ExternalIPAddr()
	AddPortMapping()
}

//搜索网关设备
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

//获得公网ip地址
func ExternalIPAddr() {
	upnpMan := new(upnp.Upnp)
	err := upnpMan.ExternalIPAddr()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("External Network IP Address:", upnpMan.GatewayOutsideIP)
	}
}

//添加一个端口映射
func AddPortMapping() {
	mapping := new(upnp.Upnp)
	if err := mapping.AddPortMapping(55789, 55789, 0, "TCP"); err == nil {
		fmt.Println("Port mapping succeeded.")
		mapping.Reclaim()
	} else {
		fmt.Println("Port mapping failed.")
	}

}
