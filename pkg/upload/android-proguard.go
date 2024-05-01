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
	BuildUuid     string      `help:"Module Build UUID" xor:"no-build-uuid,build-uuid"`
	NoBuildUuid   bool        `help:"Upload with no Build UUID" xor:"build-uuid,no-build-uuid"`
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
	noBuildUuid bool,
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
	logger log.Logger,
) error {

	var mappingFile string
	var appManifestPathExpected string
	var err error

	for _, path := range paths {
		if utils.IsDir(path) {

			mappingPath := filepath.Join(path, "app", "build", "outputs", "mapping")

			if !utils.FileExists(mappingPath) {
				return fmt.Errorf("unable to find the mapping directory in %s", path)
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
					logger.Info(fmt.Sprintf("Found app manifest at: %s", appManifestPath))
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
								logger.Info(fmt.Sprintf("Found app manifest at: %s", appManifestPath))
							}
						}
					}
				}
			}

		}

		// Check to see if we need to read the manifest file due to missing options
		if appManifestPath != "" && (apiKey == "" || applicationId == "" || buildUuid == "" || versionCode == "" || versionName == "") {

			logger.Info("Reading data from AndroidManifest.xml")
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
					logger.Info(fmt.Sprintf("Using %s as API key from AndroidManifest.xml", apiKey))
				}
			}

			if applicationId == "" {
				applicationId = manifestData.ApplicationId

				if applicationId != "" {
					logger.Info(fmt.Sprintf("Using %s as application ID from AndroidManifest.xml", applicationId))
				}
			}

			if noBuildUuid {
				buildUuid = ""
				logger.Info("No build ID will be used")
			} else if buildUuid == "" {
				for i := range manifestData.Application.MetaData.Name {
					if manifestData.Application.MetaData.Name[i] == "com.bugsnag.android.BUILD_UUID" {
						buildUuid = manifestData.Application.MetaData.Value[i]
					}
				}

				if len(dexFiles) == 0 && variant != "" {
					dexFiles = android.FindVariantDexFiles(mappingFile, variant)
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
						logger.Info(fmt.Sprintf("Using %s as build ID from classes.dex", buildUuid))
					}
				} else {
					logger.Info(fmt.Sprintf("Using %s as build UUID from AndroidManifest.xml", buildUuid))
				}
			}

			if versionCode == "" {
				versionCode = manifestData.VersionCode

				if versionCode != "" {

					logger.Info(fmt.Sprintf("Using %s as version code from AndroidManifest.xml", versionCode))
				}
			}

			if versionName == "" {
				versionName = manifestData.VersionName

				if versionName != "" {
					logger.Info(fmt.Sprintf("Using %s as version name from AndroidManifest.xml", versionName))
				}
			}
		}
		logger.Info(fmt.Sprintf("Compressing %s", mappingFile))

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

		err = server.ProcessFileRequest(endpoint+"/proguard", uploadOptions, fileFieldData, timeout, retries, outputFile, dryRun, logger)

		if err != nil {
			if strings.Contains(err.Error(), "404 Not Found") {
				logger.Info(fmt.Sprintf("Trying %s", endpoint))
				err = server.ProcessFileRequest(endpoint, uploadOptions, fileFieldData, timeout, retries, outputFile, dryRun, logger)
			}
		}

		if err != nil {
			return err
		}
	}
	return nil
}
