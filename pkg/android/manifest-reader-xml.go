package android

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
)

type Manifest struct {
	XMLName     xml.Name    `xml:"manifest"`
	Package     string      `xml:"package,attr"`
	VersionCode string      `xml:"versionCode,attr"`
	VersionName string      `xml:"versionName,attr"`
	Application Application `xml:"application"`
}

type Application struct {
	XMLName  xml.Name `xml:"application"`
	MetaData MetaData `xml:"meta-data"`
}

type MetaData struct {
	XMLName xml.Name `xml:"meta-data"`
	Name    []string `xml:"name,attr"`
	Value   []string `xml:"value,attr"`
}

// ParseAndroidManifestXML - Pulls information from a human-readable xml file into a struct
func ParseAndroidManifestXML(manifestFile string) (*Manifest, error) {
	file, err := os.Open(manifestFile)

	if err != nil {
		return nil, fmt.Errorf("unable to open " + manifestFile + " : " + err.Error())
	}

	defer file.Close()

	data, err := ioutil.ReadAll(file)

	if err != nil {
		return nil, fmt.Errorf("unable to read " + manifestFile + " : " + err.Error())
	}

	manifestData := Manifest{}

	err = xml.Unmarshal(data, &manifestData)

	if err != nil {
		return nil, fmt.Errorf("unable to parse data from " + manifestFile + " : " + err.Error())
	}

	return &manifestData, nil
}
