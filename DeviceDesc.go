package upnp

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"strings"
)

type DeviceDesc struct {
	upnp *Upnp
}

func (this *DeviceDesc) Send() bool {
	request := this.BuildRequest()
	response, _ := http.DefaultClient.Do(request)
	resultBody, _ := ioutil.ReadAll(response.Body)
	if response.StatusCode == 200 {
		this.resolve(string(resultBody))
		return true
	}
	return false
}
func (this *DeviceDesc) BuildRequest() *http.Request {
	// Request header
	header := http.Header{}
	header.Set("Accept", "text/html, image/gif, image/jpeg, *; q=.2, */*; q=.2")
	header.Set("User-Agent", "go upnp")
	header.Set("Host", this.upnp.Gateway.Host)
	header.Set("Connection", "keep-alive")

	// Request
	request, _ := http.NewRequest("GET", "http://"+this.upnp.Gateway.Host+this.upnp.Gateway.DeviceDescUrl, nil)
	request.Header = header
	// request := http.Request{Method: "GET", Proto: "HTTP/1.1",
	// 	Host: this.upnp.Gateway.Host, Url: this.upnp.Gateway.DeviceDescUrl, Header: header}
	return request
}

func (this *DeviceDesc) resolve(resultStr string) {
	inputReader := strings.NewReader(resultStr)

	// Read from the file as follows：
	// content, err := ioutil.ReadFile("studygolang.xml")
	// decoder := xml.NewDecoder(bytes.NewBuffer(content))

	lastLabel := ""

	ISUpnpServer := false

	IScontrolURL := false
	var controlURL string //`controlURL`
	// var eventSubURL string //`eventSubURL`
	// var SCPDURL string     //`SCPDURL`

	decoder := xml.NewDecoder(inputReader)
	for t, err := decoder.Token(); err == nil && !IScontrolURL; t, err = decoder.Token() {
		switch token := t.(type) {
		// Processing element start (label)
		case xml.StartElement:
			if ISUpnpServer {
				name := token.Name.Local
				lastLabel = name
			}

		// Process element end (label)
		case xml.EndElement:
			// log.Println("End tag:", token.Name.Local)
		// Processing character data (here is the text of the element)
		case xml.CharData:
			// Get url after the other tags will not be processed
			content := string([]byte(token))

			// Find the service that provides port mapping
			if content == this.upnp.Gateway.ServiceType {
				ISUpnpServer = true
				continue
			}
			//urn:upnp-org:serviceId:WANIPConnection
			if ISUpnpServer {
				switch lastLabel {
				case "controlURL":

					controlURL = content
					IScontrolURL = true
				case "eventSubURL":
					// eventSubURL = content
				case "SCPDURL":
					// SCPDURL = content
				}
			}
		default:
			// TODO:?
		}
	}
	this.upnp.CtrlUrl = controlURL
}
