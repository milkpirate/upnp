## UPnP protocol
A simple implementation of the UPnP protocol as a Golang library.  Add port mappings for NAT devices.
Look for a gateway device, check if it supports UPnP, and if so, add port mappings.

## Examples:

### 1. add a port mapping
~~~ go
mapping := new(upnp.Upnp)
if err := mapping.AddPortMapping(55789, 55789, "TCP"); err == nil {
	fmt.Println("success !")
	// remove port mapping in gateway
	mapping.Reclaim()
} else {
	fmt.Println("failed:", err.Error())
}
~~~

### 2. search gateway device.
~~~ go
upnpMan := new(upnp.Upnp)
err := upnpMan.SearchGateway()
if err != nil {
	fmt.Println(err.Error())
} else {
	fmt.Println("local ip address: ", upnpMan.LocalHost)
	fmt.Println("gateway ip address: ", upnpMan.Gateway.Host)
}
~~~
### 3. get an internet ip address in gateway.
~~~ go
upnpMan := new(upnp.Upnp)
err := upnpMan.ExternalIPAddr()
if err != nil {
	fmt.Println(err.Error())
} else {
	fmt.Println("internet ip address: ", upnpMan.GatewayOutsideIP)
}
~~~