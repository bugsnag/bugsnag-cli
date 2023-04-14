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
	ApplicationId  string            `help:"Module application identifier"`
	AndroidNdkPath string            `help:"Path to Android NDK installation ($ANDROID_NDK_ROOT)"`
	AppManifest    string            `help:"Path to app manifest file" type:"path"`
	Path           utils.UploadPaths `arg:"" name:"path" help:"Path to directory or file to upload" type:"path" default:"."`
	ProjectRoot    string            `help:"path to remove from the beginning of the filenames in the mapping file" type:"path"`
	Variant        string            `help:"Build type, like 'debug' or 'release'"`
	VersionCode    string            `help:"Module version code"`
	VersionName    string            `help:"Module version name"`
}

// ProcessAndroidNDK - Processes Android NDK symbol files
func ProcessAndroidNDK(paths []string, androidNdkPath string, appManifest string, variant string, projectRoot string, applicationId string, versionCode string, versionName string, endpoint string, timeout int, retries int, overwrite bool, apiKey string, failOnUploadError bool) error {

	uploadFileOptions := make(map[string]map[string]string)
	soFileList := make(map[string][]string)

	for _, path := range paths {
		if utils.IsDir(path) {
			// Check if we have an AndroidManifest.xml and locate it if we do not
			if appManifest == "" {
				log.Info("Locating AndroidManifest.xml in " + path)
				appManifest = filepath.Join(path, "app", "build", "intermediates", "merged_manifests")

				// Check if we have a variant and locate them if we do not
				if variant == "" {
					log.Info("Locating variant in " + appManifest)

					variants, err := android.BuildVariantsList(appManifest)

					if err != nil {
						return err
					}

					if len(variants) > 1 {
						return fmt.Errorf("more than one variant found. Please specify using `--variant`")
					}

					variant = variants[0]
				}

				appManifest = filepath.Join(appManifest, variant, "AndroidManifest.xml")
			}

			if variant == "" {
				if filepath.Base(appManifest) == "AndroidManifest.xml" {
					variantPath := filepath.Join(appManifest, "..")

					variants, err := android.BuildVariantsList(variantPath)

					if err != nil {
						return err
					}

					if len(variants) > 1 {
						return fmt.Errorf("more than one variant found. Please specify using `--variant`")
					}

					variant = variants[0]
				} else {
					fmt.Errorf("missing variant. Please specify using `--variant`")
				}
			}

			//	Build file list to process
			log.Info("Building file list for " + variant)

			mergedNativeLibs := filepath.Join(path, "app", "build", "intermediates", "merged_native_libs", variant)

			files, err := utils.BuildFileList([]string{mergedNativeLibs})

			if err != nil {
				return fmt.Errorf("error building file list for " + variant + ". " + err.Error())
			}

			var soFiles []string

			// Build two maps containing the required information to process the uploads
			for _, file := range files {
				// Ensure that we're only dealing with .so files
				if filepath.Ext(file) == ".so" && !strings.HasSuffix(file, ".sym.so") {
					uploadFileOptions[variant] = map[string]string{}
					uploadFileOptions[variant]["AndroidManifest"] = appManifest
					soFiles = append(soFiles, file)
					soFileList[variant] = soFiles
				}
			}

			// Stop processing this file and skip to the next
			continue
		} else if filepath.Ext(path) == ".so" && !strings.HasSuffix(path, ".sym.so") {
			if appManifest == "" {
				log.Info("Locating AndroidManifest.xml")
				appManifest = filepath.Join(path, "..", "..", "..", "..", "..", "..", "..", "..", "..", "app", "build", "intermediates", "merged_manifests")
			}

			// Check if we have a variant and locate them if we do not
			if variant == "" {
				if strings.Contains(path, filepath.Join("build", "intermediates", "merged_native_libs")) {
					variant = filepath.Base(filepath.Join(path, "..", "..", "..", ".."))
				} else {
					return fmt.Errorf("unable to determine variant name. Please specify using `--variant`")
				}
			}

			appManifest = filepath.Join(appManifest, variant, "AndroidManifest.xml")

			var soFiles []string

			uploadFileOptions[variant] = map[string]string{}
			uploadFileOptions[variant]["AndroidManifest"] = appManifest
			soFiles = append(soFiles, path)
			soFileList[variant] = soFiles

			// Stop processing this file and skip to the next
			continue
		} else {
			// TODO better log for when the file is not a .so file
			log.Warn("Skipping " + path + " as it is not a .so file")
			continue
		}
	}

	numberOfVariants := len(uploadFileOptions)

	if numberOfVariants < 1 {
		log.Info("Nothing to process")
		return nil
	}

	// Check NDK path is set
	androidNdkPath, err := android.GetAndroidNDKRoot(androidNdkPath)

	if err != nil {
		return err
	}

	log.Info("Android NDK Path: " + androidNdkPath)

	// Find objcopy within NDK path
	log.Info("Locating objcopy within Android NDK path")

	objCopyPath, err := android.BuildObjcopyPath(androidNdkPath)

	if err != nil {
		return err
	}

	log.Info("Objcopy Path: " + objCopyPath)

	for variant, config := range uploadFileOptions {
		log.Info("Processing files for variant: " + variant)

		// Check if the appManifest path exists
		if !utils.FileExists(appManifest) {
			return fmt.Errorf(appManifest + " does not exist on the system. Please specify using `--app-manifest`")
		}

		androidManifestData, err := android.ParseAndroidManifestXML(config["AndroidManifest"])

		if err != nil {
			return err
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

		if applicationId == "" {
			applicationId = androidManifestData.AppId
		}

		if versionCode == "" {
			versionCode = androidManifestData.VersionCode
		}

		if versionName == "" {
			versionName = androidManifestData.VersionName
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

			uploadOptions := utils.BuildAndroidNDKUploadOptions(apiKey, applicationId, versionName, versionCode, projectRoot, filepath.Base(file), overwrite)

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
