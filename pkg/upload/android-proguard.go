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
	Path          utils.Paths `arg:"" name:"path" help:"The path to the directory or file to upload" type:"path" default:"."`
	ApplicationId string      `help:"A unique application ID, usually the package name, of the application"`
	AppManifest   string      `help:"The path to a manifest file (AndroidManifest.xml) from which to obtain build information" type:"path"`
	BuildUuid     string      `help:"A unique identifier for this build of the application" xor:"no-build-uuid,build-uuid"`
	NoBuildUuid   bool        `help:"Prevents the automatically generated build UUID being uploaded with the build" xor:"build-uuid,no-build-uuid"`
	DexFiles      []string    `help:"The path to classes.dex files or directory used to calculate a build UUID" type:"path" default:""`
	Variant       string      `help:"The build type/flavor (e.g. debug, release) used to disambiguate the between built files when searching the project directory"`
	VersionCode   string      `help:"The version code of this build of the application"`
	VersionName   string      `help:"The version of the application"`
}

func ProcessAndroidProguard(
	apiKey string,
	options AndroidProguardMapping,
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

	for _, path := range options.Path {
		if utils.IsDir(path) {

			mappingPath := filepath.Join(path, "app", "build", "outputs", "mapping")

			if !utils.FileExists(mappingPath) {
				return fmt.Errorf("unable to find the mapping directory in %s", path)
			}

			if options.Variant == "" {
				options.Variant, err = android.GetVariantDirectory(mappingPath)

				if err != nil {
					return err
				}
			}

			mappingFile = filepath.Join(mappingPath, options.Variant, "mapping.txt")

			if !utils.FileExists(mappingFile) {
				return fmt.Errorf("unable to find mapping file in the specified project directory")
			}

			if options.AppManifest == "" {
				appManifestPathExpected = filepath.Join(path, "app", "build", "intermediates", "merged_manifests", options.Variant, "AndroidManifest.xml")
				if utils.FileExists(appManifestPathExpected) {
					options.AppManifest = appManifestPathExpected
					logger.Info(fmt.Sprintf("Found app manifest at: %s", options.AppManifest))
				}
			}

		} else {
			mappingFile = path

			if options.AppManifest == "" {
				if options.Variant == "" {
					//	Set the mergeNativeLibPath based off the file location e.g. outputs/mapping/<options.Variant>/mapping.txt
					mergedManifestPath := filepath.Join(path, "..", "..", "..", "..", "intermediates", "merged_manifests")

					if filepath.Base(mergedManifestPath) == "merged_manifests" {
						options.Variant, err = android.GetVariantDirectory(mergedManifestPath)
						if err == nil {
							appManifestPathExpected = filepath.Join(mergedManifestPath, options.Variant, "AndroidManifest.xml")
							if utils.FileExists(appManifestPathExpected) {
								options.AppManifest = appManifestPathExpected
								logger.Info(fmt.Sprintf("Found app manifest at: %s", options.AppManifest))
							}
						}
					}
				}
			}

		}

		// Check to see if we need to read the manifest file due to missing options
		if options.AppManifest != "" && (apiKey == "" || options.ApplicationId == "" || options.BuildUuid == "" || options.VersionCode == "" || options.VersionName == "") {

			logger.Debug("Reading data from AndroidManifest.xml")
			manifestData, err := android.ParseAndroidManifestXML(options.AppManifest)

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
					logger.Debug(fmt.Sprintf("Using %s as API key from AndroidManifest.xml", apiKey))
				}
			}

			if options.ApplicationId == "" {
				options.ApplicationId = manifestData.ApplicationId

				if options.ApplicationId != "" {
					logger.Debug(fmt.Sprintf("Using %s as application ID from AndroidManifest.xml", options.ApplicationId))
				}
			}

			if options.NoBuildUuid {
				options.BuildUuid = ""
				logger.Info("No build ID will be used")
			} else if options.BuildUuid == "" {
				for i := range manifestData.Application.MetaData.Name {
					if manifestData.Application.MetaData.Name[i] == "com.bugsnag.android.BUILD_UUID" {
						options.BuildUuid = manifestData.Application.MetaData.Value[i]
					}
				}

				if len(options.DexFiles) == 0 && options.Variant != "" {
					options.DexFiles = android.FindVariantDexFiles(mappingFile, options.Variant)
				}

				if options.BuildUuid == "" && len(options.DexFiles) > 0 {
					safeDexFile, err := android.GetDexFiles(options.DexFiles)
					if err != nil {
						return err
					}

					signature, err := android.GetAppSignatureFromFiles(safeDexFile)
					if err != nil {
						return err
					}

					options.BuildUuid = fmt.Sprintf("%x", signature)

					if options.BuildUuid != "" {
						logger.Debug(fmt.Sprintf("Using %s as build ID from classes.dex", options.BuildUuid))
					}
				} else {
					logger.Debug(fmt.Sprintf("Using %s as build UUID from AndroidManifest.xml", options.BuildUuid))
				}
			}

			if options.VersionCode == "" {
				options.VersionCode = manifestData.VersionCode

				if options.VersionCode != "" {

					logger.Debug(fmt.Sprintf("Using %s as version code from AndroidManifest.xml", options.VersionCode))
				}
			}

			if options.VersionName == "" {
				options.VersionName = manifestData.VersionName

				if options.VersionName != "" {
					logger.Debug(fmt.Sprintf("Using %s as version name from AndroidManifest.xml", options.VersionName))
				}
			}
		}
		logger.Info(fmt.Sprintf("Compressing %s", mappingFile))

		outputFile, err := utils.GzipCompress(mappingFile)

		if err != nil {
			return err
		}

		uploadOptions, err := utils.BuildAndroidProguardUploadOptions(apiKey, options.ApplicationId, options.VersionName, options.VersionCode, options.BuildUuid, overwrite)

		if err != nil {
			return err
		}

		fileFieldData := make(map[string]server.FileField)
		fileFieldData["proguard"] = server.LocalFile(outputFile)

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
