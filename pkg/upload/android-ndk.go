package upload

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/bugsnag/bugsnag-cli/pkg/android"
	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/server"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
)

type AndroidNdkMapping struct {
	AndroidNdkRoot  string            `help:"Path to Android NDK installation ($ANDROID_NDK_ROOT)"`
	AppManifestPath string            `help:"(required) Path to directory or file to upload" type:"path"`
	Configuration   string            `help:"Build type, like 'debug' or 'release'"`
	Path            utils.UploadPaths `arg:"" name:"path" help:"(required) Path to directory or file to upload" type:"path"`
	ProjectRoot     string            `help:"path to remove from the beginning of the filenames in the mapping file" type:"path"`
	VersionCode     string            `help:"Module version code"`
	VersionName     string            `help:"Module version name"`
}

// ProcessAndroidNDK - Processes Android NDK symbol files
func ProcessAndroidNDK(paths []string, androidNdkRoot string, appManifestPath string, configuration string, projectRoot string, versionCode string, versionName string, endpoint string, timeout int, retries int, overwrite bool, apiKey string, failOnUploadError bool) error {

	// Check if we have project root
	if projectRoot == "" {
		return fmt.Errorf("`--project-root` missing from options")
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

	uploadFileOptions := make(map[string]map[string]string)
	soFileList := make(map[string][]string)

	for _, path := range paths {
		if utils.IsDir(path) {
			if filepath.Base(path) == "merged_native_libs" {
				log.Info("Building variants list")

				variants, err := android.BuildVariantsList(path)

				if err != nil {
					log.Error(err.Error(), 1)
				}

				for _, variant := range variants {
					log.Info("Building file list for variant: " + variant)
					fileList, err := utils.BuildFileList([]string{filepath.Join(path, variant)})
					if err != nil {
						log.Error("error building file list for variant: "+variant, 1)
					}

					var soFiles []string

					for _, file := range fileList {
						if filepath.Ext(file) == ".so" && !strings.HasSuffix(file, ".sym.so") {
							uploadFileOptions[variant] = map[string]string{}
							uploadFileOptions[variant]["androidManifestPath"] = filepath.Join(path, "..", "merged_manifests", variant, "AndroidManifest.xml")
							uploadFileOptions[variant]["outputMetadataPath"] = filepath.Join(path, "..", "merged_manifests", variant, "output-metadata.json")
							soFiles = append(soFiles, file)
							soFileList[variant] = soFiles
						}
					}
				}

			} else {
				log.Error("unsupported folder structure provided, expected /path/to/merged_native_libs. Actual: "+path, 1)
			}
		} else if filepath.Ext(path) == ".so" {
			if configuration == "" {
				log.Warn("`--configuration` missing from options for " + path)
				log.Info("Skipping " + path)
				continue
			}
			if appManifestPath == "" {
				log.Warn("`--app-manifest-path` missing from options for " + path)
				log.Info("Skipping " + path)
				continue
			}

			uploadFileOptions[configuration] = map[string]string{}
			uploadFileOptions[configuration]["androidManifestPath"] = appManifestPath
			uploadFileOptions[configuration]["outputMetadataPath"] = filepath.Join(appManifestPath, "../output-metadata.json")
			var soFiles []string
			soFiles = append(soFiles, path)
			soFileList[configuration] = soFiles
		}
	}

	numberOfVariants := len(uploadFileOptions)

	if numberOfVariants < 1 {
		log.Info("No variants to process")
		return nil
	}

	log.Info("Processing " + strconv.Itoa(numberOfVariants) + " variant(s)")

	for variant, config := range uploadFileOptions {
		log.Info("Processing files for variant: " + variant)

		log.Info("Gathering information from AndroidManifest.xml")
		androidManifestData, err := android.ParseAndroidManifestXML(config["androidManifestPath"])

		if err != nil {
			return err
		}

		numberOfFiles := len(soFileList[variant])

		if numberOfFiles < 1 {
			log.Info("No files to process for variant: " + variant)
			continue
		}

		for _, file := range soFileList[variant] {
			log.Info("Extracting debug info from " + filepath.Base(file) + " using objcopy")
			outputFile, err := android.Objcopy(objCopyPath, file)

			if err != nil {
				log.Error("failed to process file, "+file+" using objcopy. "+err.Error(), 1)
			}

			log.Info("Uploading debug information for " + filepath.Base(file))

			uploadOptions := utils.BuildAndroidNDKUploadOptions(apiKey, androidManifestData.AppId, androidManifestData.VersionName, androidManifestData.VersionCode, projectRoot, filepath.Base(file), overwrite)

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
