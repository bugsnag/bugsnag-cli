package main

import (
	"fmt"
	"io"
)

type DartSymbol struct {
	Path []string `arg:"" name:"path" help:"Path to directory or file to upload" type:"path"`

	AppVersion string `help:"Application version"`
	AppVersionCode string `help:"Module version code (Android only)"`
	AppBundleVersion string `help:"App bundle version (Apple platforms only)"`
	Platform string `help:"platform the program was built for"`
	BuildID string `help:"foobar"`
}

func DartUpload(url string, apiKey string, buildId string, path []string) {
	var requestFieldData = map[string]string{}
	requestFieldData["buildId"] = buildId
	requestFieldData["apiKey"] = apiKey
	for _, p := range path {
		request, err := BuildFileRequest("https://upload.bugsnag.com/dart-symbol", requestFieldData, "symbolFile", p)

		if err != nil {
			break
		}

		response, err := SendRequest(request)

		b, err := io.ReadAll(response.Body)
		fmt.Println(string(b))

	}
}