package api

import (
	"encoding/xml"
	"io"
	"os"
	"log"
)

var configFilePath = "config/endpoints.xml"
var endpointFile io.Reader
var xmlEndpoints XMLEndpoints

func init(){

	loadConfigFile()
}

func loadConfigFile(){

	var err error

	endpointFile,err = os.Open(configFilePath)
	if err != nil {
		log.Fatal("could not read config endpointFile (config/endpoints.xml)")
	}

	xmlEndpoints,err = readEndpoints(endpointFile)

	if err != nil {
		log.Println("error in reading endpoints.xml config endpointFile: ")
		log.Fatal(err)
	}
		log.Printf("Loading %d endpoints:",len(xmlEndpoints.Endpoints))
	for _,endpoint := range xmlEndpoints.Endpoints{
		log.Printf("type:%s path:%s method:%s",endpoint.Type,endpoint.Path,endpoint.Method)
	}

}

func findEndpoints(path string) []XMLEndpoint{

	var list []XMLEndpoint

	for _,endpoint := range xmlEndpoints.Endpoints{
		if endpoint.Path == path{
			list = append(list,endpoint)
		}
	}
	return list
}

func provideEndpoints() XMLEndpoints{

	return xmlEndpoints

}

func readEndpoints(reader io.Reader) (XMLEndpoints, error) {
	var xmlEndpoints XMLEndpoints
	if err := xml.NewDecoder(reader).Decode(&xmlEndpoints); err != nil {
		return XMLEndpoints{}, err
	}

	return xmlEndpoints, nil
}

type XMLEndpoint struct {
	XMLName xml.Name `xml:"endpoint"`
	Type string `xml:"type,attr"`
	Path string `xml:"path,attr"`
	Method string `xml:"method,attr"`
}

type XMLEndpoints struct {
	XMLName xml.Name `xml:"endpoints"`
	Endpoints []XMLEndpoint `xml:"endpoint"`
}
