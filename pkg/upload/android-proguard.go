package upload

import (
	"fmt"
	"github.com/bugsnag/bugsnag-cli/pkg/android"
	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/server"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
	"path/filepath"
)

type AndroidProguardMapping struct {
	ApplicationId   string            `help:"Module application identifier"`
	AppManifestPath string            `help:"Path to app manifest file" type:"path"`
	MappingPath     string            `help:"Path to app mapping file"`
	BuildUuid       string            `help:"Module Build UUID"`
	Configuration   string            `help:"Build type, like 'debug' or 'release'"`
	Path            utils.UploadPaths `arg:"" name:"path" help:"Path to directory or file to upload" type:"path" default:"."`
	VersionCode     string            `help:"Module version code"`
	VersionName     string            `help:"Module version name"`
}

func ProcessAndroidProguard(paths []string, appManifestPath string, mappingPath string, buildUuid string, configuration string, appId string, versionCode string, versionName string, endpoint string, timeout int, retries int, overwrite bool, apiKey string, failOnUploadError bool, dryRun bool) error {

	uploadFileOptions := make(map[string]string)

	for _, path := range paths {

		if apiKey == "" {
			getApiKeyFromManifest = true
		}

		if appManifestPath == "" {
			log.Info("Locating Android manifest")

			if utils.FileExists(filepath.Join(path, "app", "build", "intermediates", "merged_manifests")) {
				appManifestPath = filepath.Join(path, "app", "build", "intermediates", "merged_manifests")
			} else {
				return fmt.Errorf("unable to find AndroidManifest.xml. Please specify using `--app-manifest-path` ")
			}

			if configuration == "" {
				variants, err := android.BuildVariantsList(appManifestPath)

				if err != nil {
					return err
				}

				if len(variants) > 1 {
					fmt.Println(variants)
					return fmt.Errorf("more than one variant found. Please specify using `--configuration`")
				}

				configuration = variants[0]
			}

			appManifestPath = filepath.Join(appManifestPath, configuration, "AndroidManifest.xml")
		}

		log.Info("Compressing " + config["mappingPath"])

		outputFile, err := utils.GzipCompress(config["mappingPath"])

		if err != nil {
			return err
		}

		log.Info("Uploading debug information for " + config["mappingPath"])

		uploadOptions := utils.BuildAndroidProguardUploadOptions(apiKey, androidManifestData.ApplicationId, androidManifestData.VersionName, androidManifestData.VersionCode, buildUuid, overwrite)

		fileFieldData := make(map[string]string)
		fileFieldData["proguard"] = outputFile

		requestStatus := server.ProcessRequest(endpoint, uploadOptions, fileFieldData, timeout)

		if requestStatus != nil {
			if numberOfVariants > 1 && failOnUploadError {
				return requestStatus
			} else {
				return fmt.Errorf("unable to find mapping.txt. Please specify using `--mapping-path` ")
			}
		}

		uploadFileOptions["androidManifestPath"] = appManifestPath
		uploadFileOptions["mappingPath"] = mappingPath

	}

	log.Info("Processing mapping.txt for variant: " + configuration)

	androidManifestData, err := android.ParseAndroidManifestXML(uploadFileOptions["androidManifestPath"])

	if err != nil {
		return err
	}

	if versionCode == "" {
		versionCode = androidManifestData.VersionCode
	}

	if versionName == "" {
		versionName = androidManifestData.VersionName
	}

	if appId == "" {
		appId = androidManifestData.AppId
	}

	if buildUuid == "" {
		for i := range androidManifestData.Application.MetaData.Name {
			if androidManifestData.Application.MetaData.Name[i] == "com.bugsnag.android.BUILD_UUID" {
				buildUuid = androidManifestData.Application.MetaData.Value[i]
			}
		}
	}

	log.Info("Compressing " + uploadFileOptions["mappingPath"])

	outputFile, err := utils.GzipCompress(uploadFileOptions["mappingPath"])

	if err != nil {
		return err
	}

	log.Info("Uploading debug information for " + uploadFileOptions["mappingPath"])

	uploadOptions := utils.BuildAndroidProguardUploadOptions(apiKey, appId, versionName, versionCode, buildUuid, overwrite)

	fileFieldData := make(map[string]string)
	fileFieldData["proguard"] = outputFile

	requestStatus := server.ProcessRequest(endpoint, uploadOptions, fileFieldData, timeout)

	if requestStatus != nil {
		return requestStatus
	} else {
		log.Success(uploadFileOptions["mappingPath"] + " uploaded")
	}
	return nil
}
