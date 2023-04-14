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
	ApplicationId  string            `help:"Module application identifier"`
	AndroidNdkPath string            `help:"Path to Android NDK installation ($ANDROID_NDK_ROOT)"`
	AppManifest    string            `help:"Path to app manifest file" type:"path"`
	Path           utils.UploadPaths `arg:"" name:"path" help:"Path to directory or file to upload" type:"path" default:"."`
	ProjectRoot    string            `help:"path to remove from the beginning of the filenames in the mapping file" type:"path"`
	Variant        string            `help:"Build type, like 'debug' or 'release'"`
	VersionCode    string            `help:"Module version code"`
	VersionName    string            `help:"Module version name"`
}

func ProcessAndroidNDK(apiKey string, applicationId string, androidNdkPath string, appManifest string, paths []string, projectRoot string, variant string, versionCode string, versionName string, failOnUploadError bool, endpoint string, retries int, timeout int, overwrite bool) error {

	var fileList []string
	var mergeNativeLibPath string

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

	for _, path := range paths {
		// If dir:
		if utils.IsDir(path) {
			mergeNativeLibPath = filepath.Join(path, "app", "build", "intermediates", "merged_native_libs")
			// Does merged_native_libs exist? If not, we may as well fail here...
			if !utils.FileExists(mergeNativeLibPath) {
				return fmt.Errorf("")
			}

			// If variant not provided
			if variant == "" {
				//	Get a single variant from the merged_native_libs directory, or error (if 0 or >1)
				variant, err = GetVariant(mergeNativeLibPath)

				if err != nil {
					return err
				}
			}

			// Set the upload file path(s)for the variant, or error if not exists
			fileList, err = utils.BuildFileList([]string{filepath.Join(mergeNativeLibPath, variant)})

			if err != nil {
				return fmt.Errorf("error building file list for variant: " + variant + ". " + err.Error())
			}

			// If manifest not provided
			if appManifest == "" {
				//	Get the expected path to the manifest using variant name
				appManifest = filepath.Join(path, "app", "build", "intermediates", "merged_manifests", variant, "AndroidManifest.xml")
			}

			if projectRoot == "" {
				projectRoot = path
			}
			// If file:
		} else if filepath.Ext(path) == ".so" && !strings.HasSuffix(path, ".sym.so") {
			// Set upload file path
			fileList = []string{path}

			// If manifest not provided
			if appManifest == "" {
				// If variant not provided
				if variant == "" {
					//	Iff the directory 2-levels above the file path is merged_native_libs
					mergeNativeLibPath = filepath.Join(path, "..", "..", "..", "..", "..")

					if filepath.Base(mergeNativeLibPath) == "merged_native_libs" {
						// Get variant name from directory 1-level above the file path
						variant, err = GetVariant(mergeNativeLibPath)

						if err != nil {
							return err
						}
					}
				}
				//	Get the expected path to the manifest using path relative to .so file and variant name
				appManifest = filepath.Join(mergeNativeLibPath, "..", "merged_manifests", variant, "AndroidManifest.xml")
			}

			if projectRoot == "" {
				return fmt.Errorf("project root missing")
			}
		}

		// Do we need a manifest or have all options been provided?
		if apiKey == "" {
			//   If so, check if it exists and read it
			manifestData, err := android.ParseAndroidManifestXML(appManifest)
			if err != nil {
				return err
			}

			for key, value := range manifestData.Application.MetaData.Name {
				if value == "com.bugsnag.android.API_KEY" {
					apiKey = manifestData.Application.MetaData.Value[key]
				}
			}
		}

		if applicationId == "" {
			//   If so, check if it exists and read it
			manifestData, err := android.ParseAndroidManifestXML(appManifest)
			if err != nil {
				return err
			}

			applicationId = manifestData.AppId
		}

		if versionCode == "" {
			//   If so, check if it exists and read it
			manifestData, err := android.ParseAndroidManifestXML(appManifest)
			if err != nil {
				return err
			}

			versionCode = manifestData.VersionCode
		}

		if versionName == "" {
			//   If so, check if it exists and read it
			manifestData, err := android.ParseAndroidManifestXML(appManifest)
			if err != nil {
				return err
			}

			versionName = manifestData.VersionName
		}

		numberOfFiles := len(fileList)

		if numberOfFiles < 1 {
			log.Info("No files to process")
			return nil
		}

		// Upload .so file(s)
		for _, file := range fileList {
			log.Info("Extracting debug info from " + filepath.Base(file) + " using objcopy")

			outputFile, err := android.Objcopy(objCopyPath, file)

			if err != nil {
				return fmt.Errorf("failed to process file, " + file + " using objcopy : " + err.Error())
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
