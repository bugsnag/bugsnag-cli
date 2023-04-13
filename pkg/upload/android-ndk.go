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

type AndroidNdkMapping struct {
	ApplicationId   string            `help:"Module application identifier"`
	AndroidNdkRoot  string            `help:"Path to Android NDK installation ($ANDROID_NDK_ROOT)"`
	AppManifestPath string            `help:"Path to app manifest file" type:"path"`
	Configuration   string            `help:"Build type, like 'debug' or 'release'"`
	Path            utils.UploadPaths `arg:"" name:"path" help:"Path to directory or file to upload" type:"path" default:"."`
	ProjectRoot     string            `help:"path to remove from the beginning of the filenames in the mapping file" type:"path"`
	VersionCode     string            `help:"Module version code"`
	VersionName     string            `help:"Module version name"`
}

// ProcessAndroidNDK - Processes Android NDK symbol files
func ProcessAndroidNDK(paths []string, androidNdkRoot string, appManifestPath string, configuration string, projectRoot string, appId string, versionCode string, versionName string, endpoint string, timeout int, retries int, overwrite bool, apiKey string, failOnUploadError bool) error {

	uploadFileOptions := make(map[string]string)
	var soFiles []string

	for _, path := range paths {

		if projectRoot == "" {
			projectRoot = path
		}

		androidNdkRoot, err := android.GetAndroidNDKRoot(androidNdkRoot)

		if err != nil {
			return err
		}

		log.Info("Using Android NDK located here: " + androidNdkRoot)

		log.Info("Locating objcopy within Android NDK path")

		objCopyPath, err := android.BuildObjcopyPath(androidNdkRoot)

		if err != nil {
			return err
		}

		log.Info("Using objcopy located: " + objCopyPath)

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
					return fmt.Errorf("more than one variant found. Please specify using `--configuration`")
				}

				configuration = variants[0]
			}

			appManifestPath = filepath.Join(appManifestPath, configuration, "AndroidManifest.xml")
		}

		log.Info("Building file list for variant: " + configuration)

		mergedNativeLibPath := filepath.Join(path, "app", "build", "intermediates", "merged_native_libs", configuration)

		fileList, err := utils.BuildFileList([]string{mergedNativeLibPath})

		if err != nil {
			return fmt.Errorf("error building file list for variant: " + configuration)
		}

		outputMetadataPath := filepath.Join(path, "app", "build", "merged_manifests", configuration, "output-metadata.json")

		for _, file := range fileList {
			if filepath.Ext(file) == ".so" && !strings.HasSuffix(file, ".sym.so") {
				uploadFileOptions["androidManifestPath"] = appManifestPath
				uploadFileOptions["outputMetadataPath"] = outputMetadataPath
				soFiles = append(soFiles, file)
			}
		}

		log.Info("Processing files for variant: " + configuration)

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

		if apiKey == "" {
			for key, value := range androidManifestData.Application.MetaData.Name {
				if value == "com.bugsnag.android.API_KEY" {
					apiKey = androidManifestData.Application.MetaData.Value[key]
				}
			}

			if apiKey == "" {
				return fmt.Errorf("no API key provided")
			}
		}

		numberOfFiles := len(soFiles)

		if numberOfFiles < 1 {
			log.Info("No files to process for variant: " + configuration)
			return nil
		}

		for _, file := range soFiles {
			log.Info("Extracting debug info from " + filepath.Base(file) + " using objcopy")
			outputFile, err := android.Objcopy(objCopyPath, file)

			if err != nil {
				log.Error("failed to process file, "+file+" using objcopy. "+err.Error(), 1)
			}

			log.Info("Uploading debug information for " + filepath.Base(file))

			uploadOptions := utils.BuildAndroidNDKUploadOptions(apiKey, appId, versionName, versionCode, projectRoot, filepath.Base(file), overwrite)

			fileFieldData := make(map[string]string)
			fileFieldData["soFile"] = outputFile

			requestStatus := server.ProcessRequest(endpoint, uploadOptions, fileFieldData, timeout)

			if requestStatus != nil {
				if numberOfFiles > 1 && failOnUploadError {
					return requestStatus
				} else {
					log.Warn(requestStatus.Error())
				}
			} else {
				log.Success(filepath.Base(file) + " uploaded")
			}
		}
	}

	return nil
}
