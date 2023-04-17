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
	ApplicationId string            `help:"Module application identifier"`
	AppManifest   string            `help:"Path to app manifest file" type:"path"`
	BuildUuid     string            `help:"Module Build UUID"`
	Path          utils.UploadPaths `arg:"" name:"path" help:"Path to directory or file to upload" type:"path" default:"."`
	Variant       string            `help:"Build type, like 'debug' or 'release'"`
	VersionCode   string            `help:"Module version code"`
	VersionName   string            `help:"Module version name"`
}

func ProcessAndroidProguard(apiKey string, applicationId string, appManifestPath string, buildUuid string, paths []string, variant string, versionCode string, versionName string, endpoint string, retries int, timeout int, overwrite bool, dryRun bool) error {

	var mappingFile string
	var requestStatus error
	var err error

	for _, path := range paths {
		if utils.IsDir(path) {

			mappingPath := filepath.Join(path, "app", "build", "outputs", "mapping")

			if !utils.FileExists(mappingPath) {
				return fmt.Errorf("unable to find the merged_native_libs in " + path)
			}

			if variant == "" {
				variant, err = android.GetVariant(mappingPath)

				if err != nil {
					return err
				}
			}

			mappingFile = filepath.Join(mappingPath, variant, "mapping.txt")

			if appManifestPath == "" {
				//	Get the expected path to the manifest using variant name from the given path
				appManifestPath = filepath.Join(path, "app", "build", "intermediates", "merged_manifests", variant, "AndroidManifest.xml")
			}
		} else if filepath.Base(path) == "mapping.txt" {
			mappingFile = path

			if appManifestPath == "" {
				if variant == "" {
					//	Set the mergeNativeLibPath based off the file location e.g. outputs/mapping/<variant>/mapping.txt
					mergedManifestPath := filepath.Join(path, "..", "..", "..", "..", "intermediates", "merged_manifests")

					if filepath.Base(mergedManifestPath) == "merged_manifests" {
						variant, err = android.GetVariant(mergedManifestPath)

						log.Info(mergedManifestPath)
						log.Info(variant)

						if err == nil {
							appManifestPath = filepath.Join(mergedManifestPath, variant, "AndroidManifest.xml")
						}
					}
				}
			}

		}

		// Check to see if we need to read the manifest file due to missing options
		if apiKey == "" || applicationId == "" || buildUuid == "" || versionCode == "" || versionName == "" {

			if variant == "" {
				return fmt.Errorf("missing variant. Please specify using `--variant`")
			}

			log.Info("Reading data from AndroidManifest.xml")
			manifestData, err := android.ParseAndroidManifestXML(appManifestPath)

			if err != nil {
				return err
			}

			if apiKey == "" {
				log.Info("Setting API key from AndroidManifest.xml")
				for key, value := range manifestData.Application.MetaData.Name {
					if value == "com.bugsnag.android.API_KEY" {
						apiKey = manifestData.Application.MetaData.Value[key]
					}
				}
			}

			if applicationId == "" {
				log.Info("Setting application ID from AndroidManifest.xml")
				applicationId = manifestData.ApplicationId
			}

			if buildUuid == "" {
				log.Info("Setting build UUID from AndroidManifest.xml")
				for i := range manifestData.Application.MetaData.Name {
					if manifestData.Application.MetaData.Name[i] == "com.bugsnag.android.BUILD_UUID" {
						buildUuid = manifestData.Application.MetaData.Value[i]
					}
				}
			}

			if versionCode == "" {
				log.Info("Setting version code from AndroidManifest.xml")
				versionCode = manifestData.VersionCode
			}

			if versionName == "" {
				log.Info("Setting version name from AndroidManifest.xml")
				versionName = manifestData.VersionName
			}
		}

		log.Info("Compressing " + mappingFile)

		outputFile, err := utils.GzipCompress(mappingFile)

		if err != nil {
			return err
		}

		log.Info("Uploading debug information for " + mappingFile)

		uploadOptions := utils.BuildAndroidProguardUploadOptions(apiKey, applicationId, versionName, versionCode, buildUuid, overwrite)

		fileFieldData := make(map[string]string)
		fileFieldData["proguard"] = outputFile

		if dryRun {
			requestStatus = nil
		} else {
			requestStatus = server.ProcessRequest(endpoint, uploadOptions, fileFieldData, timeout)
		}

		if requestStatus != nil {
			return requestStatus
		} else {
			log.Success(mappingFile + " uploaded")
		}
	}
	return nil
}
