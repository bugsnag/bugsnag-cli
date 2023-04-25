package upload

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bugsnag/bugsnag-cli/pkg/android"
	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/server"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
)

type AndroidNdkMapping struct {
	ApplicationId  string            `help:"Module application identifier"`
	AndroidNdkRoot string            `help:"Path to Android NDK installation ($ANDROID_NDK_ROOT)"`
	AppManifest    string            `help:"Path to app manifest file" type:"path"`
	Path           utils.UploadPaths `arg:"" name:"path" help:"Path to directory or file to upload" type:"path" default:"."`
	ProjectRoot    string            `help:"path to remove from the beginning of the filenames in the mapping file" type:"path"`
	Variant        string            `help:"Build type, like 'debug' or 'release'"`
	VersionCode    string            `help:"Module version code"`
	VersionName    string            `help:"Module version name"`
}

func ProcessAndroidNDK(apiKey string, applicationId string, androidNdkRoot string, appManifestPath string, paths []string, projectRoot string, variant string, versionCode string, versionName string, endpoint string, failOnUploadError bool, retries int, timeout int, overwrite bool, dryRun bool) error {

	var fileList []string
	var symbolFileList []string
	var mergeNativeLibPath string
	var err error
	var workingDir string
	var appManifestPathExpected string
	var objCopyPath string

	if dryRun {
		log.Info("Performing dry run - no files will be uploaded")
	}

	for _, path := range paths {
		if utils.IsDir(path) {
			mergeNativeLibPath = filepath.Join(path, "app", "build", "intermediates", "merged_native_libs")

			// Check to see if we can find app/build/intermediates/merged_native_libs from the given path
			if !utils.FileExists(mergeNativeLibPath) {
				return fmt.Errorf("unable to find the merged_native_libs in " + path)
			}

			if variant == "" {
				variant, err = android.GetVariant(mergeNativeLibPath)

				if err != nil {
					return err
				}
			}

			fileList, err = utils.BuildFileList([]string{filepath.Join(mergeNativeLibPath, variant)})

			if err != nil {
				return fmt.Errorf("error building file list for variant: " + variant + ". " + err.Error())
			}

			if appManifestPath == "" {
				appManifestPathExpected = filepath.Join(path, "app", "build", "intermediates", "merged_manifests", variant, "AndroidManifest.xml")
				if utils.FileExists(appManifestPathExpected) {
					appManifestPath = appManifestPathExpected
					log.Info("Found app manifest at: " + appManifestPath)
				}
			}

			if projectRoot == "" {
				projectRoot = path
			}

		} else {
			fileList = []string{path}

			if appManifestPath == "" {
				if variant == "" {
					//	Set the mergeNativeLibPath based off the file location e.g. merged_native_libs/<variant/out/lib/<arch>/
					mergeNativeLibPath = filepath.Join(path, "..", "..", "..", "..", "..")

					if filepath.Base(mergeNativeLibPath) == "merged_native_libs" {
						variant, err = android.GetVariant(mergeNativeLibPath)

						if err == nil {
							appManifestPathExpected = filepath.Join(mergeNativeLibPath, "..", "merged_manifests", variant, "AndroidManifest.xml")
							if utils.FileExists(appManifestPathExpected) {
								appManifestPath = appManifestPathExpected
								log.Info("Found app manifest at: " + appManifestPath)
							}
						}

						if projectRoot == "" {
							// Setting projectRoot to the suspected root of the project
							projectRoot = filepath.Join(mergeNativeLibPath, "..", "..", "..", "..")
						}
					}
				}
			}
		}

		if projectRoot != "" {
			log.Info("Using " + projectRoot + " as the project root")
		}

		// Check to see if we need to read the manifest file due to missing options
		if appManifestPath != "" && (apiKey == "" || applicationId == "" || versionCode == "" || versionName == "") {

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

				log.Info("Using " + apiKey + " as API key from AndroidManifest.xml")
			}

			if applicationId == "" {
				applicationId = manifestData.ApplicationId
				log.Info("Using " + applicationId + " as application ID from AndroidManifest.xml")
			}

			if versionCode == "" {
				versionCode = manifestData.VersionCode
				log.Info("Using " + versionCode + " as version code from AndroidManifest.xml")
			}

			if versionName == "" {
				versionName = manifestData.VersionName
				log.Info("Using " + versionName + " as version name from AndroidManifest.xml")
			}
		}

		// Process .so files through objcopy to create .sym files, filtering any other file type
		for _, file := range fileList {
			if strings.HasSuffix(file, ".so.sym") {
				symbolFileList = append(symbolFileList, file)
			} else if filepath.Ext(file) == ".so" {
				// Check NDK path is set
				if objCopyPath == "" {
					androidNdkRoot, err = android.GetAndroidNDKRoot(androidNdkRoot)

					if err != nil {
						return err
					}

					objCopyPath, err = android.BuildObjcopyPath(androidNdkRoot)

					if err != nil {
						return err
					}

					log.Info("Located objcopy within Android NDK path: " + androidNdkRoot)
				}

				log.Info("Extracting debug info from " + filepath.Base(file) + " using objcopy")

				if workingDir == "" {
					workingDir, err = os.MkdirTemp("", "bugsnag-cli-ndk-*")

					if err != nil {
						return fmt.Errorf("error creating temporary working directory " + err.Error())
					}

					defer os.RemoveAll(workingDir)
				}

				outputFile, err := android.Objcopy(objCopyPath, file, workingDir)

				if err != nil {
					return fmt.Errorf("failed to process file, " + file + " using objcopy : " + err.Error())
				}

				symbolFileList = append(symbolFileList, outputFile)
			} else {
				log.Warn("Skipping unsupported file: " + file)
			}
		}

		numberOfFiles := len(symbolFileList)

		if numberOfFiles < 1 {
			log.Info("No library files found to process")
			continue
		}

		// Upload processed .so.sym files
		for _, file := range symbolFileList {
			log.Info("Uploading debug information for " + filepath.Base(file))

			uploadOptions, err := utils.BuildAndroidNDKUploadOptions(apiKey, applicationId, versionName, versionCode, projectRoot, filepath.Base(file), overwrite)

			if err != nil {
				return err
			}

			fileFieldData := make(map[string]string)
			fileFieldData["soFile"] = file

			if dryRun {
				err = nil
			} else {
				err = server.ProcessRequest(endpoint+"/ndk-symbol", uploadOptions, fileFieldData, timeout)
			}

			if err != nil {
				if numberOfFiles > 1 && failOnUploadError {
					return err
				} else {
					log.Warn(err.Error())
				}
			} else {
				log.Success(filepath.Base(file) + " uploaded")
			}

		}
	}

	return nil
}
