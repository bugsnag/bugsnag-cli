package android

import (
	"fmt"
	"path/filepath"

	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/server"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
)

func ProcessProguard(apiKey string, variant string, outputPath string, aabManifestData map[string]string, overwrite bool, timeout int, numberOfVariants int, endpoint string, failOnUploadError bool) error {
	log.Info("Processing Proguard mapping for " + variant)

	proguardMappingPath := filepath.Join(outputPath, "BUNDLE-METADATA", "com.android.tools.build.obfuscation", "proguard.map")

	if !utils.FileExists(proguardMappingPath) {
		return fmt.Errorf(proguardMappingPath + " does not exist")
	}

	log.Info("Compressing " + proguardMappingPath)

	outputFile, err := utils.GzipCompress(proguardMappingPath)

	if err != nil {
		return err
	}

	log.Info("Uploading debug information for " + outputFile)

	uploadOptions := utils.BuildAndroidProguardUploadOptions(apiKey, aabManifestData["package"], aabManifestData["versionName"], aabManifestData["versionCode"], aabManifestData["buildUuid"], overwrite)

	fileFieldData := make(map[string]string)
	fileFieldData["proguard"] = outputFile

	requestStatus := server.ProcessRequest(endpoint, uploadOptions, fileFieldData, timeout)

	if requestStatus != nil {
		if numberOfVariants > 1 && failOnUploadError {
			return requestStatus
		} else {
			log.Warn(requestStatus.Error())
		}
	} else {
		log.Success(proguardMappingPath + " uploaded")
	}

	return nil
}
