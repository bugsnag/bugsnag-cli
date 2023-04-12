package android

import (
	"fmt"
	"path/filepath"

	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/server"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
)

func ProcessProguard(apiKey string, configuration string, outputPath string, appId string, versionCode string, versionName string, buildUuid string, mappingPath string, overwrite bool, timeout int, endpoint string, failOnUploadError bool) error {
	log.Info("Processing mapping.txt for variant: " + configuration)

	proguardMappingPath := filepath.Join(outputPath, "BUNDLE-METADATA", "com.android.tools.build.obfuscation", "proguard.map")

	if !utils.FileExists(proguardMappingPath) {
		return fmt.Errorf(proguardMappingPath + " does not exist")
	}

	log.Info("Compressing " + mappingPath)

	outputFile, err := utils.GzipCompress(mappingPath)

	if err != nil {
		return err
	}

	log.Info("Uploading debug information for " + mappingPath)

	uploadOptions := utils.BuildAndroidProguardUploadOptions(apiKey, appId, versionName, versionCode, buildUuid, overwrite)

	fileFieldData := make(map[string]string)
	fileFieldData["proguard"] = outputFile

	requestStatus := server.ProcessRequest(endpoint, uploadOptions, fileFieldData, timeout)

	if requestStatus != nil {
		return requestStatus
	} else {
		log.Success(mappingPath + " uploaded")
	}
	return nil
}
