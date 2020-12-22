package upnp

import (
	"io/ioutil"
	"log"

	"encoding/xml"
	"net/http"
	"strconv"
	"strings"

	"github.com/davecgh/go-spew/spew"
)

type GetGenericPortMappingEntry struct {
	upnp *Upnp
}

func (this *GetGenericPortMappingEntry) Send(index string) PortMappingEntry {
	request := this.buildRequest(index)
	response, _ := http.DefaultClient.Do(request)
	resultBody, _ := ioutil.ReadAll(response.Body)
	var portmap PortMappingEntry
	spew.Dump(resultBody)
	if response.StatusCode == 200 {

		portmap, err := this.resolve(string(resultBody))
		if err != nil {
			log.Println(err.Error())
			return portmap
		}
		return portmap
	}
	return portmap
}
func (this *GetGenericPortMappingEntry) buildRequest(index string) *http.Request {
	//请求头
	header := http.Header{}
	header.Set("Accept", "text/html, image/gif, image/jpeg, *; q=.2, */*; q=.2")
	header.Set("SOAPAction", `"urn:schemas-upnp-org:service:WANIPConnection:2#GetGenericPortMappingEntry"`)
	header.Set("Content-Type", "text/xml")
	header.Set("Connection", "Close")
	header.Set("Content-Length", "")

	//请求体
	body := Node{Name: "SOAP-ENV:Envelope",
		Attr: map[string]string{"xmlns:SOAP-ENV": `"http://schemas.xmlsoap.org/soap/envelope/"`,
			"SOAP-ENV:encodingStyle": `"http://schemas.xmlsoap.org/soap/encoding/"`}}
	childOne := Node{Name: `SOAP-ENV:Body`}
	childTwo := Node{Name: `m:GetGenericPortMappingEntry`,
		Attr: map[string]string{"xmlns:m": `"urn:schemas-upnp-org:service:WANIPConnection:2"`}}
	childList1 := Node{Name: "NewPortMappingIndex", Content: index}

	childTwo.AddChild(childList1)
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

func (this *GetGenericPortMappingEntry) resolve(result string) (PortMappingEntry, error) {

	res := &GenericRespEnvelope{}
	err := xml.Unmarshal([]byte(result), res)

	if err != nil {
		return res.Body.GetResponse, err
	}

	return res.Body.GetResponse, nil
}

type GenericRespEnvelope struct {
	XMLName xml.Name
	Body    GenericBody
}

type GenericBody struct {
	XMLName     xml.Name
	GetResponse PortMappingEntry `xml:"GetGenericPortMappingEntryResponse"`
}
