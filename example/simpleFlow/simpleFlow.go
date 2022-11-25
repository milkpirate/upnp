package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"
)

func chk(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	youleTest()
}

func youleTest() {
	lAddr := "192.168.1.100"
	rAddr := "192.168.1.2"
	//---------------------------------------------------------
	//      Search gateway equipment
	//---------------------------------------------------------

	searchDevice(lAddr+":9981", "239.255.255.250:1900")

	//---------------------------------------------------------
	//      View the device description
	//---------------------------------------------------------

	// readDeviceDesc(rAddr + ":1900")

	//---------------------------------------------------------
	//      Check the device status SOAPAction: "urn:schemas-upnp-org:service:WANIPConnection:1#GetStatusInfo"\r\n
	//---------------------------------------------------------
	// getDeviceStatusInfo(rAddr + ":1900")
	getDeviceStatusInfo(rAddr + ":56688")

	addPortMapping(rAddr + ":56688")

	time.Sleep(time.Second * 10)

	remotePort(rAddr + ":56688")
}

func searchDevice(localAddr, remoteAddr string) string {
	fmt.Println("Searching for gateway device...")
	msg := "M-SEARCH * HTTP/1.1\r\n" +
		"HOST: 239.255.255.250:1900\r\n" +
		"ST: urn:schemas-upnp-org:device:InternetGatewayDevice:1\r\n" +
		"MAN: \"ssdp:discover\"\r\n" +
		"MX: 3\r\n" + // seconds to delay response
		"\r\n"
	remotAddr, err := net.ResolveUDPAddr("udp", remoteAddr)
	chk(err)
	locaAddr, err := net.ResolveUDPAddr("udp", localAddr)
	chk(err)
	conn, err := net.ListenUDP("udp", locaAddr)
	chk(err)
	_, err = conn.WriteToUDP([]byte(msg), remotAddr)
	chk(err)
	buf := make([]byte, 1024)
	_, _, err = conn.ReadFromUDP(buf)
	chk(err)
	defer conn.Close()
	fmt.Println(string(buf))
	return string(buf)
}

func getDeviceStatusInfo(rAddr string) {

	fmt.Println("Fetching device status...")

	readMappingBody := `<?xml version="1.0"?>
	<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/" SOAP-ENV:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/">
	<SOAP-ENV:Body>
	<m:GetStatusInfo xmlns:m="urn:schemas-upnp-org:service:WANIPConnection:1">
	</m:GetStatusInfo></SOAP-ENV:Body></SOAP-ENV:Envelope>`

	client := &http.Client{}
	// The third parameter sets the body part
	request, _ := http.NewRequest("POST", "http://"+rAddr+"/ipc", strings.NewReader(readMappingBody))
	request.Proto = "HTTP/1.1"
	request.Host = rAddr

	request.Header.Set("Accept", "text/html, image/gif, image/jpeg, *; q=.2, */*; q=.2")
	request.Header.Set("Content-Type", "text/xml")
	request.Header.Set("SOAPAction", "\"urn:schemas-upnp-org:service:WANIPConnection:1#GetStatusInfo\"")

	request.Header.Set("Connection", "Close")
	request.Header.Set("Content-Length", fmt.Sprint(len([]byte(readMappingBody))))

	response, _ := client.Do(request)

	body, _ := ioutil.ReadAll(response.Body)
	//bodystr := string(body)
	fmt.Println(response.StatusCode)
	if response.StatusCode == 200 {
		fmt.Println(response.Header)
		fmt.Println(string(body))
	}
}

func addPortMapping(rAddr string) {

	fmt.Println("Adding a port mapping...")

	readMappingBody := `<?xml version="1.0"?>
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/" SOAP-ENV:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/">
<SOAP-ENV:Body>
<m:AddPortMapping xmlns:m="urn:schemas-upnp-org:service:WANIPConnection:1">
<NewExternalPort>6991</NewExternalPort>
<NewInternalPort>6991</NewInternalPort>
<NewProtocol>TCP</NewProtocol>
<NewEnabled>1</NewEnabled>
<NewInternalClient>192.168.1.100</NewInternalClient>
<NewLeaseDuration>0</NewLeaseDuration>
<NewPortMappingDescription>test</NewPortMappingDescription>
<NewRemoteHost></NewRemoteHost>
</m:AddPortMapping>
</SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

	client := &http.Client{}
	// The third parameter sets the body part
	request, _ := http.NewRequest("POST", "http://"+rAddr+"/ipc", strings.NewReader(readMappingBody))
	request.Proto = "HTTP/1.1"
	request.Host = rAddr

	request.Header.Set("Accept", "text/html, image/gif, image/jpeg, *; q=.2, */*; q=.2")
	request.Header.Set("Content-Type", "text/xml")
	request.Header.Set("SOAPAction", `"urn:schemas-upnp-org:service:WANIPConnection:1#AddPortMapping"`)

	request.Header.Set("Connection", "Close")
	request.Header.Set("Content-Length", fmt.Sprint(len([]byte(readMappingBody))))

	response, _ := client.Do(request)

	body, _ := ioutil.ReadAll(response.Body)
	//bodystr := string(body)
	fmt.Println(response.StatusCode)
	if response.StatusCode == 200 {
		fmt.Println(response.Header)
		fmt.Println(string(body))
	}
}

func remotePort(rAddr string) {
	fmt.Println("Deleting a port mapping...")

	readMappingBody := `<?xml version="1.0"?>
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/" SOAP-ENV:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/">
<SOAP-ENV:Body>
<m:DeletePortMapping xmlns:m="urn:schemas-upnp-org:service:WANIPConnection:1">
<NewExternalPort>6991</NewExternalPort>
<NewProtocol>TCP</NewProtocol>
<NewRemoteHost></NewRemoteHost>
</m:DeletePortMapping>
</SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

	client := &http.Client{}
	// The third parameter sets the body part
	request, _ := http.NewRequest("POST", "http://"+rAddr+"/ipc", strings.NewReader(readMappingBody))
	request.Proto = "HTTP/1.1"
	request.Host = rAddr

	request.Header.Set("Accept", "text/html, image/gif, image/jpeg, *; q=.2, */*; q=.2")
	request.Header.Set("Content-Type", "text/xml")
	request.Header.Set("SOAPAction", `"urn:schemas-upnp-org:service:WANIPConnection:1#DeletePortMapping"`)

	request.Header.Set("Connection", "Close")
	request.Header.Set("Content-Length", fmt.Sprint(len([]byte(readMappingBody))))

	response, _ := client.Do(request)

	body, _ := ioutil.ReadAll(response.Body)
	//bodystr := string(body)
	fmt.Println(response.StatusCode)
	if response.StatusCode == 200 {
		fmt.Println(response.Header)
		fmt.Println(string(body))
	}
}
