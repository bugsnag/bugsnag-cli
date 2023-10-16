package upload

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/bugsnag/bugsnag-cli/pkg/android"
	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/server"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
)

type AndroidProguardMapping struct {
	ApplicationId string      `help:"Module application identifier"`
	AppManifest   string      `help:"Path to app manifest file" type:"path"`
	BuildUuid     string      `help:"Module Build UUID"`
	DexFiles      []string    `help:"Path to classes.dex files or directory" type:"path" default:""`
	Path          utils.Paths `arg:"" name:"path" help:"Path to directory or file to upload" type:"path" default:"."`
	Variant       string      `help:"Build type, like 'debug' or 'release'"`
	VersionCode   string      `help:"Module version code"`
	VersionName   string      `help:"Module version name"`
}

func ProcessAndroidProguard(
	apiKey string,
	applicationId string,
	appManifestPath string,
	buildUuid string,
	dexFiles []string,
	paths []string,
	variant string,
	versionCode string,
	versionName string,
	endpoint string,
	retries int,
	timeout int,
	overwrite bool,
	dryRun bool,
) error {

	var mappingFile string
	var appManifestPathExpected string
	var err error

	for _, path := range paths {
		if utils.IsDir(path) {

			mappingPath := filepath.Join(path, "app", "build", "outputs", "mapping")

			if !utils.FileExists(mappingPath) {
				return fmt.Errorf("unable to find the mapping directory in " + path)
			}

			if variant == "" {
				variant, err = android.GetVariantDirectory(mappingPath)

				if err != nil {
					return err
				}
			}

			mappingFile = filepath.Join(mappingPath, variant, "mapping.txt")

			if !utils.FileExists(mappingFile) {
				return fmt.Errorf("unable to find mapping file in the specified project directory")
			}

			if appManifestPath == "" {
				appManifestPathExpected = filepath.Join(path, "app", "build", "intermediates", "merged_manifests", variant, "AndroidManifest.xml")
				if utils.FileExists(appManifestPathExpected) {
					appManifestPath = appManifestPathExpected
					log.Info("Found app manifest at: " + appManifestPath)
				}
			}

		} else {
			mappingFile = path

			if appManifestPath == "" {
				if variant == "" {
					//	Set the mergeNativeLibPath based off the file location e.g. outputs/mapping/<variant>/mapping.txt
					mergedManifestPath := filepath.Join(path, "..", "..", "..", "..", "intermediates", "merged_manifests")

					if filepath.Base(mergedManifestPath) == "merged_manifests" {
						variant, err = android.GetVariantDirectory(mergedManifestPath)
						if err == nil {
							appManifestPathExpected = filepath.Join(mergedManifestPath, variant, "AndroidManifest.xml")
							if utils.FileExists(appManifestPathExpected) {
								appManifestPath = appManifestPathExpected
								log.Info("Found app manifest at: " + appManifestPath)
							}
						}
					}
				}
			}

		}

		// Check to see if we need to read the manifest file due to missing options
		if appManifestPath != "" && (apiKey == "" || applicationId == "" || buildUuid == "" || versionCode == "" || versionName == "") {

			log.Info("Reading data from AndroidManifest.xml")
			manifestData, err := android.ParseAndroidManifestXML(appManifestPath)

			if err != nil {
				return err
			}

			if apiKey == "" {
				for key, value := range manifestData.Application.MetaData.Name {
					if value == "com.bugsnag.android.API_KEY" {
						apiKey = manifestData.Application.MetaData.Value[key]
					}
				}

				if apiKey != "" {
					log.Info("Using " + apiKey + " as API key from AndroidManifest.xml")
				}
			}

			if applicationId == "" {
				applicationId = manifestData.ApplicationId

				if applicationId != "" {
					log.Info("Using " + applicationId + " as application ID from AndroidManifest.xml")
				}
			}

			if buildUuid == "" {
				for i := range manifestData.Application.MetaData.Name {
					if manifestData.Application.MetaData.Name[i] == "com.bugsnag.android.BUILD_UUID" {
						buildUuid = manifestData.Application.MetaData.Value[i]
					}
				}

				if buildUuid == "" && len(dexFiles) > 0 {
					safeDexFile, err := android.GetDexFiles(dexFiles)
					if err != nil {
						return err
					}

					signature, err := android.GetAppSignatureFromFiles(safeDexFile)
					if err != nil {
						return err
					}

					buildUuid = fmt.Sprintf("%x", signature)

					if buildUuid != "" {
						log.Info("Using " + buildUuid + " as build ID from classes.dex")
					}
				} else {
					log.Info("Using " + buildUuid + " as build UUID from AndroidManifest.xml")
				}
			}

			if versionCode == "" {
				versionCode = manifestData.VersionCode

				if versionCode != "" {
					log.Info("Using " + versionCode + " as version code from AndroidManifest.xml")
				}
			}

			if versionName == "" {
				versionName = manifestData.VersionName

				if versionName != "" {
					log.Info("Using " + versionName + " as version name from AndroidManifest.xml")
				}
			}
		}

		log.Info("Compressing " + mappingFile)

		outputFile, err := utils.GzipCompress(mappingFile)

		if err != nil {
			return err
		}

		uploadOptions, err := utils.BuildAndroidProguardUploadOptions(apiKey, applicationId, versionName, versionCode, buildUuid, overwrite)

		if err != nil {
			return err
		}

		fileFieldData := make(map[string]string)
		fileFieldData["proguard"] = outputFile

		err = server.ProcessRequest(endpoint+"/proguard", uploadOptions, fileFieldData, timeout, outputFile, dryRun)

		if err != nil {
			if strings.Contains(err.Error(), "404 Not Found") {
				log.Info("Trying " + endpoint)
				err = server.ProcessRequest(endpoint, uploadOptions, fileFieldData, timeout, outputFile, dryRun)
			}
		}

		if err != nil {
			return err
		} else {
			log.Success("Uploaded " + filepath.Base(mappingFile))
		}
	}
	return nil
}
