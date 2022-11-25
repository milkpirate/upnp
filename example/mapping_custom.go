package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/fdurand/upnp"
)

var mapping = new(upnp.Upnp)
var reader = bufio.NewReader(os.Stdin)

var localPort = 1990
var remotePort = 1990

func init() {

}

func main() {
	Start()
}

func Start() {
	if !CheckNet() {
		fmt.Println("Your router does not support the UPnP protocol.")
		return
	}
	fmt.Println("Local IP Address: ", mapping.LocalHost)

	ExternalIPAddr()

tag:
	if !GetInput() {
		goto tag
	}
	if !AddPortMapping(localPort, remotePort) {
		goto tag
	}

	fmt.Println("--------------------------------------")
	fmt.Println("1.  stop    stop the program and reclaim mapped port")
	fmt.Println("2.  add     add a port mapping")
	fmt.Println("3.  del     manually delete a port mapping")
	fmt.Println("\n NOTE: This progrma maps tcp ports. If you need to")
	fmt.Println("       map a UDP port, please visit：")
	fmt.Println("       http://github.com/prestonTao/upnp")
	fmt.Println("--------------------------------------")

	running := true
	for running {
		data, _, _ := reader.ReadLine()
		commands := strings.Split(string(data), " ")
		switch commands[0] {
		case "help":

		case "stop":
			running = false
			mapping.Reclaim()
		case "add":
			goto tag
		case "del":
		tagDel:
			if !GetInput() {
				goto tagDel
			}
			DelPortMapping(localPort, remotePort)
		case "cdp":
		case "dump":
		}
	}

}

/*
	检查网络是否支持upnp协议
*/
func CheckNet() bool {
	err := mapping.SearchGateway()
	if err != nil {
		return false
	} else {
		return true
	}
}

//获得公网ip地址
func ExternalIPAddr() {
	err := mapping.ExternalIPAddr()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("External network IP address: ", mapping.GatewayOutsideIP)
	}
}

/*
	得到用户输入的端口
*/
func GetInput() bool {
	var err error
	fmt.Println("Enter the local port to map: ")
	data, _, _ := reader.ReadLine()
	localPort, err = strconv.Atoi(string(data))
	if err != nil {
		fmt.Println("Invalid port, please specify a valid port between 0-65535")
		return false
	}
	if localPort < 0 || localPort > 65535 {
		fmt.Println("Invalid port, please specify a valid port between 0-65535")
		return false
	}

	fmt.Println("Please enter the external port to be mapped:")
	data, _, _ = reader.ReadLine()
	remotePort, err = strconv.Atoi(string(data))
	if err != nil {
		fmt.Println("Invalid port, please specify a valid port between 0-65535")
		return false
	}
	if remotePort < 0 || remotePort > 65535 {
		fmt.Println("Invalid port, please specify a valid port between 0-65535")
		return false
	}
	return true
}

/*
	添加一个端口映射
*/
func AddPortMapping(localPort, remotePort int) bool {
	//添加一个端口映射 localPort, remotePort, duration int, protocol string, desc string
	if err := mapping.AddPortMapping(localPort, remotePort, 10, "TCP", "wireguard"); err == nil {
		fmt.Println("Port mapped successfully")
		return true
	} else {
		fmt.Println("Port failed to map")
		return false
	}
}

/*
	删除一个端口映射
*/
func DelPortMapping(localPort, remotePort int) {
	mapping.DelPortMapping(remotePort, "TCP")
}
