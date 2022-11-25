package upnp

import (
	"errors"
	"fmt"
	"sync"

	"github.com/davecgh/go-spew/spew"
)

// Get the Gateway

// MappingPortStruct Manage all Ports
type MappingPortStruct struct {
	lock         *sync.Mutex
	mappingPorts map[string][][]int
}

// Add a port mapping record
// only map management
func (m *MappingPortStruct) addMapping(localPort, remotePort int, protocol string) {

	m.lock.Lock()
	defer m.lock.Unlock()
	if m.mappingPorts == nil {
		one := make([]int, 0)
		one = append(one, localPort)
		two := make([]int, 0)
		two = append(two, remotePort)
		portMapping := [][]int{one, two}
		m.mappingPorts = map[string][][]int{protocol: portMapping}
		return
	}
	portMapping := m.mappingPorts[protocol]
	if portMapping == nil {
		one := make([]int, 0)
		one = append(one, localPort)
		two := make([]int, 0)
		two = append(two, remotePort)
		m.mappingPorts[protocol] = [][]int{one, two}
		return
	}
	one := portMapping[0]
	two := portMapping[1]
	one = append(one, localPort)
	two = append(two, remotePort)
	m.mappingPorts[protocol] = [][]int{one, two}
}

// Delete a mapping record
// only map management
func (m *MappingPortStruct) delMapping(remotePort int, protocol string) {
	m.lock.Lock()
	defer m.lock.Unlock()
	if m.mappingPorts == nil {
		return
	}
	tmp := MappingPortStruct{lock: new(sync.Mutex)}
	mappings := m.mappingPorts[protocol]
	for i := 0; i < len(mappings[0]); i++ {
		if mappings[1][i] == remotePort {
			// the map to delete
			break
		}
		tmp.addMapping(mappings[0][i], mappings[1][i], protocol)
	}
	m.mappingPorts = tmp.mappingPorts
}
func (m *MappingPortStruct) GetAllMapping() map[string][][]int {
	return m.mappingPorts
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
func (u *Upnp) SearchGateway() (err error) {
	defer func() {
		if errTemp := recover(); errTemp != nil {
			fmt.Println("SearchGateway:", errTemp)
			err = errTemp.(error)
		}
	}()

	if u.LocalHost == "" {
		u.MappingPort = MappingPortStruct{
			lock: new(sync.Mutex),
			// mappingPorts: map[string][][]int{},
		}
		u.LocalHost = GetLocalIntenetIp()
	}
	searchGateway := SearchGateway{upnp: u}
	if searchGateway.Send() {
		return nil
	}
	return errors.New("no gateway device found")
}

// View the device description and get the control request URL
func (u *Upnp) deviceDesc() (err error) {
	if u.GatewayInsideIP == "" {
		if err := u.SearchGateway(); err != nil {
			return err
		}
	}
	device := DeviceDesc{upnp: u}
	device.Send()
	u.Active = true

	return
}

// ExternalIPAddr gets our external IP address
func (u *Upnp) ExternalIPAddr() (err error) {
	if u.CtrlUrl == "" {
		if err := u.deviceDesc(); err != nil {
			return err
		}
	}
	eia := ExternalIPAddress{upnp: u}
	eia.Send()

	return nil
}

// AddPortMapping adds a port mapping
// TODO: accept an IP address to port forward to another LAN host(internalClient)
func (u *Upnp) AddPortMapping(localPort, remotePort, duration int, internalClient string, protocol string, desc string) (err error) {
	defer func() {
		if errTemp := recover(); errTemp != nil {
			fmt.Println("AddPortMapping:", errTemp)
			err = errTemp.(error)
		}
	}()
	if u.GatewayOutsideIP == "" {
		if err := u.ExternalIPAddr(); err != nil {
			return err
		}
	}
	addPort := AddPortMapping{upnp: u}
	if isSuccess := addPort.Send(localPort, remotePort, duration, internalClient, protocol, desc); isSuccess {
		u.MappingPort.addMapping(localPort, remotePort, protocol)

		return nil
	} else {
		u.Active = false
		// fmt.Println("failed to add port mapping")
		// TODO: is it possible to get an error from gateway instead of showing our own?
		return errors.New("adding a port mapping failed")
	}
}

// DelPortMapping probably deletes a port mapping
func (u *Upnp) DelPortMapping(remotePort int, protocol string) bool {
	delMapping := DelPortMapping{upnp: u}
	isSuccess := delMapping.Send(remotePort, protocol)
	if isSuccess {
		u.MappingPort.delMapping(remotePort, protocol)
		fmt.Println("Removed a port mapping: remote:", remotePort)
	}
	return isSuccess
}

func (u *Upnp) GetListOfPortMappings(protocol string) []PortMappingEntry {
	spew.Dump(u)
	listPort := GetListOfPortMappings{upnp: u}
	portMap := listPort.Send(protocol)
	return portMap
}

func (u *Upnp) GetGenericPortMappingEntry(index string) PortMappingEntry {
	listPort := GetGenericPortMappingEntry{upnp: u}
	portMap := listPort.Send(index)
	return portMap
}

// Reclaim recycles (deletes?) a port
func (u *Upnp) Reclaim() {
	mappings := u.MappingPort.GetAllMapping()
	tcpMapping, ok := mappings["TCP"]
	if ok {
		for i := 0; i < len(tcpMapping[0]); i++ {
			u.DelPortMapping(tcpMapping[1][i], "TCP")
		}
	}
	udpMapping, ok := mappings["UDP"]
	if ok {
		for i := 0; i < len(udpMapping[0]); i++ {
			u.DelPortMapping(udpMapping[0][i], "UDP")
		}
	}
}

// GetAllMapping returns all active mappings
// TODO: for this host only?
func (u *Upnp) GetAllMapping() map[string][][]int {
	return u.MappingPort.GetAllMapping()
}
