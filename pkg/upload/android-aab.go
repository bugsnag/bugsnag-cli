package upload

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/bugsnag/bugsnag-cli/pkg/android"
	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/options"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
)

// ProcessAndroidAab processes Android AAB files for upload.
//
// It supports resolving AAB files either directly by path or by searching expected
// build output directories within the project. The function extracts metadata and
// symbol files from the AAB, then triggers uploads for NDK (.so) files and Proguard
// mapping files if found.
//
// Parameters:
//   - globalOptions: CLI options containing upload paths and settings.
//   - logger: logger instance for logging debug and info messages.
//
// Returns:
//   - error: non-nil if any step in processing or uploading fails.
func ProcessAndroidAab(globalOptions options.CLI, logger log.Logger) error {
	var manifestData map[string]string
	var aabDir string
	var aabFile string
	var err error
	aabOptions := globalOptions.Upload.AndroidAab

	for _, path := range aabOptions.Path {
		// If the path is a directory, check if it contains extracted AAB metadata or try to find the AAB file.
		if utils.IsDir(path) {
			if utils.FileExists(filepath.Join(path, "BUNDLE-METADATA")) {
				aabDir = path
			} else {
				// Search common AAB build output paths for the AAB file.
				arr := []string{"*", "build", "outputs", "bundle", "release", "*-release*.aab"}
				aabFile, err = android.FindAabPath(arr, path)
				if err != nil {
					return err
				}
			}
		} else if filepath.Ext(path) == ".aab" {
			// If path is directly an AAB file, use it.
			aabFile = path
		}

		// If we have an AAB file and no extracted directory, extract it now.
		if aabFile != "" && aabDir == "" {
			logger.Debug(fmt.Sprintf("Extracting AAB file: %s", aabFile))
			aabDir, err = utils.ExtractFile(aabFile, "aab")
			defer os.RemoveAll(aabDir) // Clean up extracted files afterward
			if err != nil {
				return err
			}
		}
	}

	// Merge upload options with metadata extracted from the AAB manifest.
	manifestData, err = android.MergeUploadOptionsFromAabManifest(
		aabDir,
		globalOptions.ApiKey,
		aabOptions.ApplicationId,
		aabOptions.BuildUuid,
		aabOptions.NoBuildUuid,
		aabOptions.VersionCode,
		aabOptions.VersionName,
		logger,
	)
	if err != nil {
		return err
	}

	// Process NDK (.so) files if present.
	soFilePath := filepath.Join(aabDir, "BUNDLE-METADATA", "com.android.tools.build.debugsymbols")
	if utils.FileExists(soFilePath) {
		logger.Debug(fmt.Sprintf("Found NDK (.so) files at: %s", soFilePath))
		soFileList, err := utils.BuildFileList([]string{soFilePath})
		if err != nil {
			return err
		}
		if len(soFileList) > 0 {
			globalOptions.Upload.AndroidNdk = options.AndroidNdkMapping{
				ApplicationId: manifestData["applicationId"],
				Path:          soFileList,
				ProjectRoot:   aabOptions.ProjectRoot,
				VersionCode:   manifestData["versionCode"],
				VersionName:   manifestData["versionName"],
			}
			globalOptions.ApiKey = manifestData["apiKey"]
			err = ProcessAndroidNDK(globalOptions, logger)
			if err != nil {
				return err
			}
		} else {
			logger.Info("No NDK (.so) files detected for upload.")
		}
	} else {
		logger.Info("No NDK (.so) files detected for upload.")
	}

	// Process Proguard mapping file if present.
	mappingFilePath := filepath.Join(aabDir, "BUNDLE-METADATA", "com.android.tools.build.obfuscation", "proguard.map")
	if utils.FileExists(mappingFilePath) {
		logger.Debug(fmt.Sprintf("Found Proguard (mapping.txt) file at: %s", mappingFilePath))
		globalOptions.Upload.AndroidProguard = options.AndroidProguardMapping{
			ApplicationId: manifestData["applicationId"],
			BuildUuid:     manifestData["buildUuid"],
			NoBuildUuid:   aabOptions.NoBuildUuid,
			DexFiles:      []string{filepath.Join(aabDir, "base", "dex")},
			Path:          []string{mappingFilePath},
			VersionCode:   manifestData["versionCode"],
			VersionName:   manifestData["versionName"],
		}
		globalOptions.ApiKey = manifestData["apiKey"]
		err = ProcessAndroidProguard(globalOptions, logger)
		if err != nil {
			return err
		}
	} else {
		logger.Info("No Proguard (mapping.txt) file detected for upload.")
	}

	return nil
}
