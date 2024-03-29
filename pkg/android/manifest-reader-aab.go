package android

import (
	"github.com/bugsnag/bugsnag-cli/pkg/proto_messages"
	"google.golang.org/protobuf/proto"
	"os"
)

func ReadAabManifest(path string) (map[string]string, error) {
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
		if data.ResourceId == 16843291 || data.Name == "versionCode" {
			aabManifestData["versionCode"] = data.GetValue()
			continue
		}

		if data.Name == "package" {
			aabManifestData["applicationId"] = data.GetValue()
			continue
		}

		if data.ResourceId == 16843292 || data.Name == "versionName" {
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

	return aabManifestData, nil
}
