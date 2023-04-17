package upload

import (
	"fmt"
	"github.com/bugsnag/bugsnag-cli/pkg/android"
	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/server"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
	"io/ioutil"
	"path/filepath"
	"strings"
)

type AndroidNdkMapping struct {
	Package        string            `help:"Module application identifier"`
	AndroidNdkRoot string            `help:"Path to Android NDK installation ($ANDROID_NDK_ROOT)"`
	AppManifest    string            `help:"Path to app manifest file" type:"path"`
	Path           utils.UploadPaths `arg:"" name:"path" help:"Path to directory or file to upload" type:"path" default:"."`
	ProjectRoot    string            `help:"path to remove from the beginning of the filenames in the mapping file" type:"path"`
	Variant        string            `help:"Build type, like 'debug' or 'release'"`
	VersionCode    string            `help:"Module version code"`
	VersionName    string            `help:"Module version name"`
}

func ProcessAndroidNDK(apiKey string, _package string, androidNdkRoot string, appManifest string, paths []string, projectRoot string, variant string, versionCode string, versionName string, endpoint string, failOnUploadError bool, retries int, timeout int, overwrite bool, dryRun bool) error {

	var fileList []string
	var mergeNativeLibPath string

	if dryRun {
		log.Info("Performing dry run")
	}

	// Check NDK path is set
	androidNdkRoot, err := android.GetAndroidNDKRoot(androidNdkRoot)

	if err != nil {
		return err
	}

	log.Info("Android NDK Path: " + androidNdkRoot)

	// Find objcopy within NDK path
	log.Info("Locating objcopy within Android NDK path")

	objCopyPath, err := android.BuildObjcopyPath(androidNdkRoot)

	if err != nil {
		return err
	}

	log.Info("Objcopy Path: " + objCopyPath)

	for _, path := range paths {
		if utils.IsDir(path) {
			mergeNativeLibPath = filepath.Join(path, "app", "build", "intermediates", "merged_native_libs")

			// Check to see if we can find app/build/intermediates/merged_native_libs from the given path
			if !utils.FileExists(mergeNativeLibPath) {
				return fmt.Errorf("unable to find the merged_native_libs in " + path)
			}

			if variant == "" {
				variant, err = GetVariant(mergeNativeLibPath)

				if err != nil {
					return err
				}
			}

			fileList, err = utils.BuildFileList([]string{filepath.Join(mergeNativeLibPath, variant)})

			if err != nil {
				return fmt.Errorf("error building file list for variant: " + variant + ". " + err.Error())
			}

			if appManifest == "" {
				//	Get the expected path to the manifest using variant name from the given path
				appManifest = filepath.Join(path, "app", "build", "intermediates", "merged_manifests", variant, "AndroidManifest.xml")
			}

			if projectRoot == "" {
				projectRoot = path
			}

		} else if filepath.Ext(path) == ".so" && !strings.HasSuffix(path, ".sym.so") {
			fileList = []string{path}

			if appManifest == "" {
				if variant == "" {
					//	Set the mergeNativeLibPath based off the file location e.g. merged_native_libs/<variant/out/lib/<arch>/
					mergeNativeLibPath = filepath.Join(path, "..", "..", "..", "..", "..")

					if filepath.Base(mergeNativeLibPath) == "merged_native_libs" {
						variant, err = GetVariant(mergeNativeLibPath)

						if err != nil {
							return err
						}
					}
				}
				appManifest = filepath.Join(mergeNativeLibPath, "..", "merged_manifests", variant, "AndroidManifest.xml")

				if utils.FileExists(mergeNativeLibPath) {
					if projectRoot == "" {
						// Setting projectRoot to the suspected root of the project
						projectRoot = filepath.Join(mergeNativeLibPath, "..", "..", "..", "..")
					}
				}
			}
		}

		log.Info("Using " + projectRoot + " as the project root")

		// Check to see if we need to read the manifest file due to missing options
		if apiKey == "" || _package == "" || versionCode == "" || versionName == "" {

			log.Info("Reading data from AndroidManifest.xml")
			manifestData, err := android.ParseAndroidManifestXML(appManifest)

			if err != nil {
				return err
			}

			if apiKey == "" {
				log.Info("Setting API key from AndroidManifest.xml")
				for key, value := range manifestData.Application.MetaData.Name {
					if value == "com.bugsnag.android.API_KEY" {
						apiKey = manifestData.Application.MetaData.Value[key]
					}
				}
			}

			if _package == "" {
				log.Info("Setting application ID from AndroidManifest.xml")
				_package = manifestData.Package
			}

			if versionCode == "" {
				log.Info("Setting version code from AndroidManifest.xml")
				versionCode = manifestData.VersionCode
			}

			if versionName == "" {
				log.Info("Setting version name from AndroidManifest.xml")
				versionName = manifestData.VersionName
			}
		}

		// Upload .so file(s)
		for _, file := range fileList {
			if filepath.Ext(path) == ".so" && !strings.HasSuffix(path, ".sym.so") {

				numberOfFiles := len(fileList)

				if numberOfFiles < 1 {
					log.Info("No files found to process")
					continue
				}

				log.Info("Extracting debug info from " + filepath.Base(file) + " using objcopy")

				outputFile, err := android.Objcopy(objCopyPath, file)

				if err != nil {
					return fmt.Errorf("failed to process file, " + file + " using objcopy : " + err.Error())
				}

				if !dryRun {
					log.Info("Uploading debug information for " + filepath.Base(file))

					uploadOptions, err := utils.BuildAndroidNDKUploadOptions(apiKey, _package, versionName, versionCode, projectRoot, filepath.Base(file), overwrite)

					if err != nil {
						return err
					}

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
		}
	}

	return nil
}

func GetVariant(path string) (string, error) {
	var variants []string

	fileInfo, err := ioutil.ReadDir(path)

	if err != nil {
		return "", err
	}

	for _, file := range fileInfo {
		variants = append(variants, file.Name())
	}

	if len(variants) > 1 {
		return "", fmt.Errorf("too many variants")
	} else if len(variants) < 1 {
		return "", fmt.Errorf("no variants")
	}

	variant := variants[0]

	if !utils.FileExists(filepath.Join(path, variant)) {
		return "", fmt.Errorf("path doesn't exist " + filepath.Join(path, variant))
	}

	return variant, nil
}
