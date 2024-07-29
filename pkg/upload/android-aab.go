package upload

import (
	"os"
	"path/filepath"

	"github.com/bugsnag/bugsnag-cli/pkg/android"
	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
)

type AndroidAabMapping struct {
	Path          utils.Paths `arg:"" name:"path" help:"The path to the AAB file to upload (or directory containing it)" type:"path" default:"."`
	ApplicationId string      `help:"A unique application ID, usually the package name, of the application"`
	BuildUuid     string      `help:"A unique identifier for this build of the application" xor:"no-build-uuid,build-uuid"`
	NoBuildUuid   bool        `help:"Prevents the automatically generated build UUID being uploaded with the build" xor:"build-uuid,no-build-uuid"`
	ProjectRoot   string      `help:"The path to strip from the beginning of source file names referenced in stacktraces on the BugSnag dashboard" type:"path"`
	VersionCode   string      `help:"The version code of this build of the application"`
	VersionName   string      `help:"The version of the application"`
}

func ProcessAndroidAab(
	apiKey string,
	options AndroidAabMapping,
	endpoint string,
	retries int,
	timeout int,
	overwrite bool,
	dryRun bool,
	logger log.Logger,
) error {

	var manifestData map[string]string
	var aabDir string
	var aabFile string
	var err error

	for _, path := range options.Path {
		// Look for AAB file if the upload command was run somewhere within the project root
		// based on an expected path of ${dir}/build/outputs/bundle/release/${dir}-release.aab
		// or ${dir}/build/outputs/bundle/release/${dir}-release-dexguard.aab
		if utils.IsDir(path) {
			if utils.FileExists(filepath.Join(path, "BUNDLE-METADATA")) {
				aabDir = path
			} else {
				arr := []string{"*", "build", "outputs", "bundle", "release", "*-release*.aab"}
				aabFile, err = android.FindAabPath(arr, path)

				if err != nil {
					return err
				}
			}
		} else if filepath.Ext(path) == ".aab" {
			aabFile = path
		}

		if aabFile != "" && aabDir == "" {
			aabDir, err = utils.ExtractFile(aabFile, "aab")

			defer os.RemoveAll(aabDir)

			if err != nil {
				return err
			}
		}
	}

	manifestData, err = android.MergeUploadOptionsFromAabManifest(aabDir, apiKey, options.ApplicationId, options.BuildUuid, options.NoBuildUuid, options.VersionCode, options.VersionName, logger)

	if err != nil {
		return err
	}

	soFilePath := filepath.Join(aabDir, "BUNDLE-METADATA", "com.android.tools.build.debugsymbols")

	if utils.FileExists(soFilePath) {
		soFileList, err := utils.BuildFileList([]string{soFilePath})

		if err != nil {
			return err
		}

		if len(soFileList) > 0 {
			err = ProcessAndroidNDK(
				manifestData["apiKey"],
				AndroidNdkMapping{
					ApplicationId: manifestData["applicationId"],
					Path:          soFileList,
					ProjectRoot:   options.ProjectRoot,
					VersionCode:   manifestData["versionCode"],
					VersionName:   manifestData["versionName"],
				},
				endpoint,
				retries,
				timeout,
				overwrite,
				dryRun,
				logger,
			)

			if err != nil {
				return err
			}
		} else {
			logger.Info("No NDK (.so) files detected for upload.")
		}
	} else {
		logger.Info("No NDK (.so) files detected for upload.")
	}

	mappingFilePath := filepath.Join(aabDir, "BUNDLE-METADATA", "com.android.tools.build.obfuscation", "proguard.map")

	if utils.FileExists(mappingFilePath) {
		err = ProcessAndroidProguard(
			manifestData["apiKey"],
			AndroidProguardMapping{
				ApplicationId: manifestData["applicationId"],
				BuildUuid:     manifestData["buildUuid"],
				NoBuildUuid:   options.NoBuildUuid,
				DexFiles:      []string{filepath.Join(aabDir, "base", "dex")},
				Path:          []string{mappingFilePath},
				VersionCode:   manifestData["versionCode"],
				VersionName:   manifestData["versionName"],
			},
			endpoint,
			retries,
			timeout,
			overwrite,
			dryRun,
			logger,
		)

		if err != nil {
			return err
		}
	} else {
		logger.Info("No Proguard (mapping.txt) file detected for upload.")
	}

	return nil
}
