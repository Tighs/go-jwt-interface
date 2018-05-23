package api

import (
	"encoding/xml"
	"io"
	"os"
	"log"
	"errors"
)

var configFilePath = "config/config.xml"
var configFile io.Reader
var config xmlConfiguration

func init(){

	loadConfigFile()
}

func loadConfigFile(){

	var err error

	configFile,err = os.Open(configFilePath)
	if err != nil {
		log.Fatal("could not read config file (config/config.xml)")
	}

	var errorList []error

	config,errorList = readConfigFile(configFile)

	if len(errorList) > 0{

		log.Println("Error occured while reading the config file:")
		for _,msg := range errorList{
			log.Println(msg)
		}
		log.Fatal("Abort due to previous errors")
	}

		log.Printf("Loading %d endpoints:",len(config.Endpoints.EndpointList))
	for _,endpoint := range config.Endpoints.EndpointList {
		log.Printf("type:%s path:%s method:%s",endpoint.Type,endpoint.Path,endpoint.Method)
	}

	log.Println("Loading Key:")
	log.Printf("KeyPair %s loaded ", config.Key.Name)
}


func findEndpoints(path string) []xmlEndpoint {

	var list []xmlEndpoint

	for _,endpoint := range config.Endpoints.EndpointList {
		if endpoint.Path == path{
			list = append(list,endpoint)
		}
	}
	return list
}

func provideEndpoints() xmlEndpoints {

	return config.Endpoints

}

func readConfigFile(reader io.Reader) (xmlConfiguration,[]error){
	var config xmlConfiguration
	var errorList []error
	if err := xml.NewDecoder(reader).Decode(&config); err != nil {
		return xmlConfiguration{},[]error{err}
	}

	for _,endpoint := range config.Endpoints.EndpointList {
		if  endpoint.Type == "" || endpoint.Path == "" ||  endpoint.Method == "" {
			errorList =append(errorList, errors.New("endpoint values must be present"))
		}
	}

	if config.Key.Name == "" || config.Key.Method == "" || config.Key.Path == ""{
		errorList =append(errorList, errors.New("key values name,method and path must be present"))
	}

	if config.Token.Expiration == 0{
		config.Token.Expiration = 30
	}

	if len(errorList) > 0{
		return xmlConfiguration{},errorList
	}
	return config,nil
}

type xmlEndpoint struct {
	XMLName xml.Name `xml:"endpoint"`
	Type string `xml:"type,attr"`
	Path string `xml:"path,attr"`
	Method string `xml:"method,attr"`
}

type xmlEndpoints struct {
	XMLName      xml.Name      `xml:"endpoints"`
	EndpointList []xmlEndpoint `xml:"endpoint"`
}

type xmlKey struct {
	XMLName xml.Name `xml:"key"`
	Name string `xml:"name,attr"`
	Method string `xml:"method,attr"`
	Size string `xml:"size,attr"`
	Path string `xml:"path,attr"`
	PassPhrase string `xml:"passphrase,attr"`
}

type xmlToken struct{
	XMLName xml.Name `xml:"token"`
	Expiration int `xml:"expiration,attr"`
	Refreshable bool `xml:"refreshable,attr"`
}

type xmlConfiguration struct {
	XMLName xml.Name `xml:"config"`
	Key xmlKey `xml:"key"`
	Endpoints xmlEndpoints `xml:"endpoints"`
	Token xmlToken `xml:"token"`
}
