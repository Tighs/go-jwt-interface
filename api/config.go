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
var endpoints xmlEndpoints
var key xmlKey

func init(){

	loadConfigFile()
}

func loadConfigFile(){

	var err error

	endpointFile,err = os.Open(configFilePath)
	if err != nil {
		log.Fatal("could not read config endpointFile (config/endpoints.xml)")
	}

	endpoints,err = readEndpoints(endpointFile)

	if err != nil {
		log.Println("error in reading endpoints.xml config endpointFile: ")
		log.Fatal(err)
	}
		log.Printf("Loading %d endpoints:",len(endpoints.Endpoints))
	for _,endpoint := range endpoints.Endpoints{
		log.Printf("type:%s path:%s method:%s",endpoint.Type,endpoint.Path,endpoint.Method)
	}

	keyFile,err = os.Open(keyFilePath)

	if err != nil {
		log.Fatal("could not read config keyFile (config/key.xml)")
	}

	key,err = readKey(keyFile)

	if err != nil {
		log.Println("error in reading key.xml: ")
		log.Fatal(err)

	}
	log.Println("Loading Key:")
	log.Printf("KeyPair %s loaded ", key.Name)
}

func findEndpoints(path string) []xmlEndpoint {

	var list []xmlEndpoint

	for _,endpoint := range endpoints.Endpoints{
		if endpoint.Path == path{
			list = append(list,endpoint)
		}
	}
	return list
}

func provideEndpoints() xmlEndpoints {

	return endpoints

}

func readEndpoints(reader io.Reader) (xmlEndpoints, error) {
	var endpoints xmlEndpoints
	if err := xml.NewDecoder(reader).Decode(&endpoints); err != nil {
		return xmlEndpoints{}, err
	}

	for _,endpoint := range endpoints.Endpoints{
		if  endpoint.Type == "" || endpoint.Path == "" ||  endpoint.Method == "" {
			return xmlEndpoints{},errors.New("endpoint values must be present")
		}
	}

	return endpoints, nil
}

func readKey(reader io.Reader) (xmlKey,error){
	var key xmlKey

	if err := xml.NewDecoder(reader).Decode(&key); err != nil {
		return xmlKey{},err
	}

	if key.Name == "" || key.Method == "" || key.Path == ""{
		return xmlKey{},errors.New("key values name,method,path must be present")
	}

	return key,nil
}

type xmlEndpoint struct {
	XMLName xml.Name `xml:"endpoint"`
	Type string `xml:"type,attr"`
	Path string `xml:"path,attr"`
	Method string `xml:"method,attr"`
}

type xmlEndpoints struct {
	XMLName xml.Name        `xml:"endpoints"`
	Endpoints []xmlEndpoint `xml:"endpoint"`
}

type xmlKey struct {
	XMLName xml.Name `xml:"key"`
	Name string `xml:"name,attr"`
	Method string `xml:"method,attr"`
	Size string `xml:"size,attr"`
	Path string `xml:"path,attr"`
	PassPhrase string `xml:"passphrase,attr"`
}
