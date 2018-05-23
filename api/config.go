package api

import (
	"encoding/xml"
	"io"
	"os"
	"log"
	"errors"
)

var configFilePath = "config/endpoints.xml"
var keyFilePath = "config/key.xml"
var endpointFile io.Reader
var keyFile io.Reader
var xmlEndpoints XMLEndpoints
var xmlKey XMLKey

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

	keyFile,err = os.Open(keyFilePath)

	if err != nil {
		log.Fatal("could not read config keyFile (config/key.xml)")
	}

	xmlKey,err = readKey(keyFile)

	if err != nil {
		log.Println("error in reading key.xml: ")
		log.Fatal(err)

	}
	log.Println("Loading Key:")
	log.Printf("KeyPair %s loaded ", xmlKey.Name)
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

	for _,endpoint := range xmlEndpoints.Endpoints{
		if  endpoint.Type == "" || endpoint.Path == "" ||  endpoint.Method == "" {
			return XMLEndpoints{},errors.New("endpoint values must be present")
		}
	}

	return xmlEndpoints, nil
}

func readKey(reader io.Reader) (XMLKey,error){
	var xmlKey XMLKey

	if err := xml.NewDecoder(reader).Decode(&xmlKey); err != nil {
		return XMLKey{},err
	}

	if xmlKey.Name == "" || xmlKey.Method == "" || xmlKey.Path == ""{
		return XMLKey{},errors.New("key values name,method,path must be present")
	}

	return xmlKey,nil
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

type XMLKey struct {
	XMLName xml.Name `xml:"key"`
	Name string `xml:"name,attr"`
	Method string `xml:"method,attr"`
	Size string `xml:"size,attr"`
	Path string `xml:"path,attr"`
	PassPhrase string `xml:"passphrase,attr"`
}
