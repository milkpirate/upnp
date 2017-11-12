package upnp

import (
	"errors"
	"fmt"
	"sync"
)

// Get the Gateway

// Manage all Ports
type MappingPortStruct struct {
	lock         *sync.Mutex
	mappingPorts map[string][][]int
}

// Add a port mapping record
// only map management
func (this *MappingPortStruct) addMapping(localPort, remotePort int, protocol string) {

	this.lock.Lock()
	defer this.lock.Unlock()
	if this.mappingPorts == nil {
		one := make([]int, 0)
		one = append(one, localPort)
		two := make([]int, 0)
		two = append(two, remotePort)
		portMapping := [][]int{one, two}
		this.mappingPorts = map[string][][]int{protocol: portMapping}
		return
	}
	portMapping := this.mappingPorts[protocol]
	if portMapping == nil {
		one := make([]int, 0)
		one = append(one, localPort)
		two := make([]int, 0)
		two = append(two, remotePort)
		this.mappingPorts[protocol] = [][]int{one, two}
		return
	}
	one := portMapping[0]
	two := portMapping[1]
	one = append(one, localPort)
	two = append(two, remotePort)
	this.mappingPorts[protocol] = [][]int{one, two}
}

// Delete a mapping record
// only map management
func (this *MappingPortStruct) delMapping(remotePort int, protocol string) {
	this.lock.Lock()
	defer this.lock.Unlock()
	if this.mappingPorts == nil {
		return
	}
	tmp := MappingPortStruct{lock: new(sync.Mutex)}
	mappings := this.mappingPorts[protocol]
	for i := 0; i < len(mappings[0]); i++ {
		if mappings[1][i] == remotePort {
			// the map to delete
			break
		}
		tmp.addMapping(mappings[0][i], mappings[1][i], protocol)
	}
	this.mappingPorts = tmp.mappingPorts
}
func (this *MappingPortStruct) GetAllMapping() map[string][][]int {
	return this.mappingPorts
}

type Upnp struct {
	Active              bool // is this upnp protocol available?
	DurationUnsupported bool
	LocalHost           string            // local (our) IP Address
	GatewayInsideIP     string            // LAN Gateway IP
	GatewayOutsideIP    string            // Gateway Public IP
	OutsideMappingPort  map[string]int    // Map the external port
	InsideMappingPort   map[string]int    // Map local port
	Gateway             *Gateway          // Gateway information
	CtrlUrl             string            // Control request URL
	MappingPort         MappingPortStruct // Existing mappings, e.g {"TCP":[1990],"UDP":[1991]}
}

// SearchGateway gets the Gateway's LAN IP address
func (this *Upnp) SearchGateway() (err error) {
	defer func(err error) {
		if errTemp := recover(); errTemp != nil {
			fmt.Println("SearchGateway:", errTemp)
			err = errTemp.(error)
		}
	}(err)

	if this.LocalHost == "" {
		this.MappingPort = MappingPortStruct{
			lock: new(sync.Mutex),
			// mappingPorts: map[string][][]int{},
		}
		this.LocalHost = GetLocalIntenetIp()
	}
	searchGateway := SearchGateway{upnp: this}
	if searchGateway.Send() {
		return nil
	}
	return errors.New("No gateway device found")
}

func (this *Upnp) deviceStatus() {

}

// View the device description and get the control request URL
func (this *Upnp) deviceDesc() (err error) {
	if this.GatewayInsideIP == "" {
		if err := this.SearchGateway(); err != nil {
			return err
		}
	}
	device := DeviceDesc{upnp: this}
	device.Send()
	this.Active = true

	return
}

// ExternalIPAddr gets our external IP address
func (this *Upnp) ExternalIPAddr() (err error) {
	if this.CtrlUrl == "" {
		if err := this.deviceDesc(); err != nil {
			return err
		}
	}
	eia := ExternalIPAddress{upnp: this}
	eia.Send()

	return nil
}

// AddPortMapping adds a port mapping
// TODO: accept an IP address to port forward to another LAN host(internalClient)
func (this *Upnp) AddPortMapping(localPort, remotePort, duration int, internalClient string, protocol string, desc string) (err error) {
	defer func(err error) {
		if errTemp := recover(); errTemp != nil {
			fmt.Println("AddPortMapping:", errTemp)
			err = errTemp.(error)
		}
	}(err)
	if this.GatewayOutsideIP == "" {
		if err := this.ExternalIPAddr(); err != nil {
			return err
		}
	}
	addPort := AddPortMapping{upnp: this}
	if issuccess := addPort.Send(localPort, remotePort, duration, internalClient, protocol, desc); issuccess {
		this.MappingPort.addMapping(localPort, remotePort, protocol)

		return nil
	} else {
		this.Active = false
		// fmt.Println("failed to add port mapping")
		// TODO: is it possible to get an error from gateway instead of showing our own?
		return errors.New("Adding a port mapping failed")
	}
}

// DelPortMapping probably deletes a port mapping
func (this *Upnp) DelPortMapping(remotePort int, protocol string) bool {
	delMapping := DelPortMapping{upnp: this}
	issuccess := delMapping.Send(remotePort, protocol)
	if issuccess {
		this.MappingPort.delMapping(remotePort, protocol)
		fmt.Println("Removed a port mapping: remote:", remotePort)
	}
	return issuccess
}

// Reclaim recycles (deletes?) a port
func (this *Upnp) Reclaim() {
	mappings := this.MappingPort.GetAllMapping()
	tcpMapping, ok := mappings["TCP"]
	if ok {
		for i := 0; i < len(tcpMapping[0]); i++ {
			this.DelPortMapping(tcpMapping[1][i], "TCP")
		}
	}
	udpMapping, ok := mappings["UDP"]
	if ok {
		for i := 0; i < len(udpMapping[0]); i++ {
			this.DelPortMapping(udpMapping[0][i], "UDP")
		}
	}
}

// GetAllMapping returns all active mappings
// TODO: for this host only?
func (this *Upnp) GetAllMapping() map[string][][]int {
	return this.MappingPort.GetAllMapping()
}
