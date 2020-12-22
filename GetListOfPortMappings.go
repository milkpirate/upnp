package upnp

import (
	"io/ioutil"
	"log"

	"encoding/xml"
	"net/http"
	"strconv"
	"strings"
)

type GetListOfPortMappings struct {
	upnp *Upnp
}

func (this *GetListOfPortMappings) Send(protocol string) []PortMappingEntry {
	request := this.buildRequest(protocol)
	response, _ := http.DefaultClient.Do(request)
	resultBody, _ := ioutil.ReadAll(response.Body)
	if response.StatusCode == 200 {
		portmap, err := this.resolve(string(resultBody))
		if err != nil {
			log.Println(err.Error())
			return nil
		}
		return portmap
	}
	return nil
}

func (this *GetListOfPortMappings) buildRequest(protocol string) *http.Request {
	//请求头
	header := http.Header{}
	header.Set("Accept", "text/html, image/gif, image/jpeg, *; q=.2, */*; q=.2")
	header.Set("SOAPAction", `"urn:schemas-upnp-org:service:WANIPConnection:2#GetListOfPortMappings"`)
	header.Set("Content-Type", "text/xml")
	header.Set("Connection", "Close")
	header.Set("Content-Length", "")

	//请求体
	body := Node{Name: "SOAP-ENV:Envelope",
		Attr: map[string]string{"xmlns:SOAP-ENV": `"http://schemas.xmlsoap.org/soap/envelope/"`,
			"SOAP-ENV:encodingStyle": `"http://schemas.xmlsoap.org/soap/encoding/"`}}
	childOne := Node{Name: `SOAP-ENV:Body`}
	childTwo := Node{Name: `m:GetListOfPortMappings`,
		Attr: map[string]string{"xmlns:m": `"urn:schemas-upnp-org:service:WANIPConnection:2"`}}
	childList1 := Node{Name: "NewStartPort", Content: "1"}
	childList2 := Node{Name: "NewEndPort", Content: "65535"}
	childList3 := Node{Name: "NewProtocol", Content: protocol}
	childList4 := Node{Name: "NewManage", Content: "1"}
	childList5 := Node{Name: "NewNumberOfPorts", Content: "65535"}

	childTwo.AddChild(childList1)
	childTwo.AddChild(childList2)
	childTwo.AddChild(childList3)
	childTwo.AddChild(childList4)
	childTwo.AddChild(childList5)
	childOne.AddChild(childTwo)
	body.AddChild(childOne)
	bodyStr := body.BuildXML()

	//请求
	request, _ := http.NewRequest("POST", "http://"+this.upnp.Gateway.Host+this.upnp.CtrlUrl,
		strings.NewReader(bodyStr))
	request.Header = header
	request.Header.Set("Content-Length", strconv.Itoa(len([]byte(bodyStr))))
	return request
}

func (this *GetListOfPortMappings) resolve(result string) ([]PortMappingEntry, error) {

	res := &MyRespEnvelope{}
	err := xml.Unmarshal([]byte(result), res)

	if err != nil {
		return nil, err
	}
	portmap := &PortMapList{}
	err = xml.Unmarshal([]byte(res.Body.GetResponse.NewPortListing), portmap)
	if err != nil {
		return nil, err
	}

	return portmap.PortMappingEntry, nil
}

type MyRespEnvelope struct {
	XMLName xml.Name
	Body    Body
}

type Body struct {
	XMLName     xml.Name
	GetResponse NewPortListing `xml:"GetListOfPortMappingsResponse"`
}

type NewPortListing struct {
	XMLName        xml.Name `xml:"GetListOfPortMappingsResponse"`
	NewPortListing string   `xml:"NewPortListing"`
}

type PortMapList struct {
	XMLName          xml.Name
	PortMappingEntry []PortMappingEntry `xml: "PortMappingEntry`
}

type PortMappingEntry struct {
	NewRemoteHost     string `xml: "NewRemoteHost"`
	NewExternalPort   string `xml: "NewExternalPort"`
	NewProtocol       string `xml: "NewProtocol"`
	NewInternalPort   string `xml: "NewInternalPort"`
	NewInternalClient string `xml: "NewInternalClient"`
	NewEnabled        string `xml: "NewEnabled"`
	NewDescription    string `xml: "NewDescription"`
	NewLeaseTime      string `xml: "NewLeaseTime"`
}
