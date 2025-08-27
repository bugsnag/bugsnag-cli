package upload

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/bugsnag/bugsnag-cli/pkg/android"
	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/options"
	"github.com/bugsnag/bugsnag-cli/pkg/server"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
)

// ProcessAndroidProguard processes and uploads Android Proguard mapping files.
//
// This function locates the Proguard mapping file(s) from given paths or directories,
// extracts metadata from AndroidManifest.xml if needed, compresses the mapping file,
// and uploads it.
//
// Parameters:
//   - options: CLI options containing upload configuration and metadata.
//   - logger: Logger instance for debug/info/error output.
//
// Returns:
//   - error: non-nil if any step fails during processing or uploading.
func ProcessAndroidProguard(options options.CLI, logger log.Logger) error {
	proguardOptions := options.Upload.AndroidProguard

	var mappingFile string
	var appManifestPathExpected string
	var err error

	for _, path := range proguardOptions.Path {
		if utils.IsDir(path) {
			// Expect mapping files under <path>/app/build/outputs/mapping
			mappingPath := filepath.Join(path, "app", "build", "outputs", "mapping")

			if !utils.FileExists(mappingPath) {
				return fmt.Errorf("unable to find the mapping directory in %s", path)
			}

			// Determine build variant if not specified
			if proguardOptions.Variant == "" {
				proguardOptions.Variant, err = android.GetVariantDirectory(mappingPath)
				if err != nil {
					return err
				}
			}

			// Compose full path to mapping.txt for the variant
			mappingFile = filepath.Join(mappingPath, proguardOptions.Variant, "mapping.txt")

			if !utils.FileExists(mappingFile) {
				return fmt.Errorf("unable to find mapping file in the specified project directory")
			}

			// Attempt to locate AndroidManifest.xml for the variant if not set
			if proguardOptions.AppManifest == "" {
				appManifestPathExpected = filepath.Join(path, "app", "build", "intermediates", "merged_manifests", proguardOptions.Variant, "AndroidManifest.xml")
				if utils.FileExists(appManifestPathExpected) {
					proguardOptions.AppManifest = appManifestPathExpected
					logger.Info(fmt.Sprintf("Found app manifest at: %s", proguardOptions.AppManifest))
				}
			}

		} else {
			// If a file path is specified directly, use it as the mapping file
			mappingFile = path

			// Try to find manifest and variant if missing
			if proguardOptions.AppManifest == "" && proguardOptions.Variant == "" {
				mergedManifestPath := filepath.Join(path, "..", "..", "..", "..", "intermediates", "merged_manifests")
				if filepath.Base(mergedManifestPath) == "merged_manifests" {
					proguardOptions.Variant, err = android.GetVariantDirectory(mergedManifestPath)
					if err == nil {
						appManifestPathExpected = filepath.Join(mergedManifestPath, proguardOptions.Variant, "AndroidManifest.xml")
						if utils.FileExists(appManifestPathExpected) {
							proguardOptions.AppManifest = appManifestPathExpected
							logger.Info(fmt.Sprintf("Found app manifest at: %s", proguardOptions.AppManifest))
						}
					}
				}
			}
		}

		// Read manifest for missing metadata: API key, application ID, build UUID, version code/name
		if proguardOptions.AppManifest != "" && (options.ApiKey == "" || proguardOptions.ApplicationId == "" || proguardOptions.BuildUuid == "" || proguardOptions.VersionCode == "" || proguardOptions.VersionName == "") {
			logger.Debug("Reading data from AndroidManifest.xml")
			manifestData, err := android.ParseAndroidManifestXML(proguardOptions.AppManifest)
			if err != nil {
				return err
			}

			if options.ApiKey == "" {
				for key, value := range manifestData.Application.MetaData.Name {
					if value == "com.bugsnag.android.API_KEY" {
						options.ApiKey = manifestData.Application.MetaData.Value[key]
					}
				}
				if options.ApiKey != "" {
					logger.Debug(fmt.Sprintf("Using %s as API key from AndroidManifest.xml", options.ApiKey))
				}
			}

			if proguardOptions.ApplicationId == "" {
				proguardOptions.ApplicationId = manifestData.ApplicationId
				if proguardOptions.ApplicationId != "" {
					logger.Debug(fmt.Sprintf("Using %s as application ID from AndroidManifest.xml", proguardOptions.ApplicationId))
				}
			}

			if proguardOptions.NoBuildUuid {
				proguardOptions.BuildUuid = ""
				logger.Info("No build ID will be used")
			} else if proguardOptions.BuildUuid == "" {
				for i := range manifestData.Application.MetaData.Name {
					if manifestData.Application.MetaData.Name[i] == "com.bugsnag.android.BUILD_UUID" {
						proguardOptions.BuildUuid = manifestData.Application.MetaData.Value[i]
					}
				}

				if len(proguardOptions.DexFiles) == 0 && proguardOptions.Variant != "" {
					proguardOptions.DexFiles = android.FindVariantDexFiles(mappingFile, proguardOptions.Variant)
				}

				if proguardOptions.BuildUuid == "" && len(proguardOptions.DexFiles) > 0 {
					safeDexFile, err := android.GetDexFiles(proguardOptions.DexFiles)
					if err != nil {
						return err
					}

					signature, err := android.GetAppSignatureFromFiles(safeDexFile)
					if err != nil {
						return err
					}

					proguardOptions.BuildUuid = fmt.Sprintf("%x", signature)

					if proguardOptions.BuildUuid != "" {
						logger.Debug(fmt.Sprintf("Using %s as build ID from classes.dex", proguardOptions.BuildUuid))
					}
				} else {
					logger.Debug(fmt.Sprintf("Using %s as build UUID from AndroidManifest.xml", proguardOptions.BuildUuid))
				}
			}

			if proguardOptions.VersionCode == "" {
				proguardOptions.VersionCode = manifestData.VersionCode
				if proguardOptions.VersionCode != "" {
					logger.Debug(fmt.Sprintf("Using %s as version code from AndroidManifest.xml", proguardOptions.VersionCode))
				}
			}

			if proguardOptions.VersionName == "" {
				proguardOptions.VersionName = manifestData.VersionName
				if proguardOptions.VersionName != "" {
					logger.Debug(fmt.Sprintf("Using %s as version name from AndroidManifest.xml", proguardOptions.VersionName))
				}
			}
		}

		logger.Info(fmt.Sprintf("Compressing %s", mappingFile))

		// Compress mapping file with gzip
		outputFile, err := utils.GzipCompress(mappingFile)
		if err != nil {
			return err
		}

		uploadOptions, err := utils.BuildAndroidProguardUploadOptions(proguardOptions.ApplicationId, proguardOptions.VersionName, proguardOptions.VersionCode, proguardOptions.BuildUuid, proguardOptions.Overwrite)

		if err != nil {
			return err
		}

		// Prepare file upload data
		fileFieldData := make(map[string]server.FileField)
		fileFieldData["proguard"] = server.LocalFile(outputFile)

		// Attempt upload to Bugsnag API
		err = server.ProcessFileRequest(options.ApiKey, "/proguard", uploadOptions, fileFieldData, outputFile, options, logger)

		// Retry at base endpoint if 404 received
		if err != nil && strings.Contains(err.Error(), "404 Not Found") {
			logger.Debug("Retrying upload for proguard at base endpoint")
			err = server.ProcessFileRequest(options.ApiKey, "", uploadOptions, fileFieldData, outputFile, options, logger)
		}

		if err != nil {
			return err
		}
	}

	return nil
}
