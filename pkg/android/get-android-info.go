package android

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"fmt"
	"github.com/bugsnag/bugsnag-cli/pkg/proto_messages"
	"google.golang.org/protobuf/proto"
	"io"
	"os"
	"path/filepath"
)

// Android Manifest reference attriubte IDs
// https://developer.android.com/reference/android/R.attr#versionCode
const AndroidVersionCodeId uint32 = 16843291

// https://developer.android.com/reference/android/R.attr#versionName
const AndroidVersionNameId uint32 = 16843292

type AndroidManifestData struct {
	XMLName       xml.Name                       `xml:"manifest"`
	ApplicationId string                         `xml:"package,attr"`
	VersionCode   string                         `xml:"versionCode,attr"`
	VersionName   string                         `xml:"versionName,attr"`
	Application   AndroidManifestApplicationData `xml:"application"`
}

type AndroidManifestApplicationData struct {
	XMLName  xml.Name                `xml:"application"`
	MetaData AndroidManifestMetaData `xml:"meta-data"`
}

type AndroidManifestMetaData struct {
	XMLName xml.Name `xml:"meta-data"`
	Name    []string `xml:"name,attr"`
	Value   []string `xml:"value,attr"`
}

func isXMLContent(buffer []byte) bool {
	xmlHeader := []byte("<manifest xmlns")
	return bytes.Contains(buffer, xmlHeader)
}

func isProtobufContent(path string) (bool, error) {
	// Open the file
	file, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer file.Close()

	// Read the first 8 bytes of the file
	buffer := make([]byte, 8)
	_, err = file.Read(buffer)
	if err != nil {
		return false, err
	}

	// Check for the protobuf identifier
	protobufIdentifier := []byte{0x0A} // Example: Protobuf identifier often starts with 0x0A
	return bytes.HasPrefix(buffer, protobufIdentifier), nil
}

// Function to extract the android manifest file from a given aab file
func GetAndroidManifestFileFromAAB(path string) (string, error) {
	aabManifestPath := filepath.Join("base", "manifest", "AndroidManifest.xml")
	outputPath := filepath.Join(path, "..", filepath.Base(aabManifestPath))

	if filepath.Ext(path) == ".aab" {
		zipData, err := zip.OpenReader(path)

		if err != nil {
			return "", err
		}

		defer zipData.Close()

		for _, file := range zipData.File {
			if file.Name == aabManifestPath {
				zippedFile, err := file.Open()
				if err != nil {
					return "", nil
				}
				defer zippedFile.Close()

				destinationFile, err := os.Create(outputPath)

				if err != nil {
					return "", err
				}

				defer destinationFile.Close()

				_, err = io.Copy(destinationFile, zippedFile)
				if err != nil {
					return "", err
				}

				return outputPath, nil
			}
		}
	}
	return "", fmt.Errorf("no aab file")
}

// getAndroidXMLData - Pulls information from a human-readable xml file into a struct
func getAndroidXMLData(manifestFile string) (*AndroidManifestData, error) {
	data, err := os.ReadFile(manifestFile)
	if err != nil {
		return nil, fmt.Errorf("unable to read %s : %w", manifestFile, err.Error())
	}

	var manifestData *AndroidManifestData
	err = xml.Unmarshal(data, &manifestData)
	if err != nil {
		return nil, fmt.Errorf("unable to parse data from %s : %w", manifestFile, err.Error())
	}

	return manifestData, nil
}

func getAndroidProtobufData(path string) (*AndroidManifestData, error) {
	aabManifestData := make(map[string]string)

	content, err := os.ReadFile(path)

	if err != nil {
		return nil, err
	}

	rawAabManifestData := &proto_messages.XmlNode{}

	err = proto.Unmarshal(content, rawAabManifestData)

	if err != nil {
		return nil, err
	}

	for _, data := range rawAabManifestData.GetElement().GetAttribute() {
		if data.ResourceId == AndroidVersionCodeId || data.Name == "versionCode" {
			aabManifestData["versionCode"] = data.GetValue()
			continue
		}

		if data.Name == "package" {
			aabManifestData["applicationId"] = data.GetValue()
			continue
		}

		if data.ResourceId == AndroidVersionNameId || data.Name == "versionName" {
			aabManifestData["versionName"] = data.GetValue()
			continue
		}
	}

	for _, level1 := range rawAabManifestData.GetElement().GetChild() {
		for _, level2 := range level1.GetElement().GetChild() {
			if level2.GetElement().GetName() == "meta-data" {
				if level2.GetElement().GetAttribute()[0].Value == "com.bugsnag.android.API_KEY" {
					aabManifestData["apiKey"] = level2.GetElement().GetAttribute()[1].Value
					continue
				}
				if level2.GetElement().GetAttribute()[0].Value == "com.bugsnag.android.BUILD_UUID" {
					aabManifestData["buildUuid"] = level2.GetElement().GetAttribute()[1].Value
					continue
				}
			}
		}
	}

	manifestData := &AndroidManifestData{
		XMLName: xml.Name{
			Space: "",
			Local: "manifest",
		},
		ApplicationId: aabManifestData["applicationId"],
		VersionCode:   aabManifestData["versionCode"],
		VersionName:   aabManifestData["versionName"],
		Application: AndroidManifestApplicationData{
			XMLName: xml.Name{
				Space: "",
				Local: "application",
			},
			MetaData: AndroidManifestMetaData{
				XMLName: xml.Name{
					Space: "",
					Local: "meta-data",
				},
				Name:  []string{"com.bugsnag.android.API_KEY", "com.bugsnag.android.BUILD_UUID"},
				Value: []string{aabManifestData["apiKey"], aabManifestData["buildUuid"]},
			},
		},
	}

	return manifestData, nil
}

func BuildAndroidInfo(path string) (*AndroidManifestData, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	buffer := make([]byte, 512)
	_, err = file.Read(buffer)

	if err != nil && err != io.EOF {
		return nil, err
	}

	contentType := isXMLContent(buffer)

	if contentType {
		androidData, err := getAndroidXMLData(path)

		if err != nil {
			return nil, err
		}

		return androidData, nil
	} else {
		contentType, err := isProtobufContent(path)

		if err != nil {
			return nil, err
		}

		if contentType {
			androidData, err := getAndroidProtobufData(path)

			return androidData, err
		}
	}

	return nil, fmt.Errorf("Unsupported file type")
}
