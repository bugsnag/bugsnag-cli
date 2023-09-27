package upload

import (
	"fmt"
	"github.com/bugsnag/bugsnag-cli/pkg/android"
	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/server"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
	"os"
	"path/filepath"
	"strings"
)

type UnityAndroidOptions struct {
	AabPath       string            `help:"Path to Android AAB file"`
	ApplicationId string            `help:"Module application identifier"`
	Arch          string            `help:"The architecture of the shared object that the symbols are for (e.g. x86, armeabi-v7a)."`
	Path          utils.UploadPaths `arg:"" name:"path" help:"(required) Path to directory or file to upload" type:"path"`
	ProjectRoot   string            `help:"path to remove from the beginning of the filenames in the mapping file" type:"path"`
	VersionCode   string            `help:"The version code for the application (Android only)."`
	VersionName   string            `help:"The version name of the application."`
	BuildUuid     string            `help:"The version name of the application."`
}

func ProcessUnityAndroid(apiKey string, aabPath string, applicationId string, versionCode string, buildUuid string, arch string, versionName string, projectRoot string, paths []string, endpoint string, failOnUploadError bool, retries int, timeout int, overwrite bool, dryRun bool) error {
	var err error
	var zipPath string
	var archList []string
	var symbolFileList []string
	var manifestData map[string]string
	var aabManifestPath string

	for _, path := range paths {
		if utils.IsDir(path) {
			zipPath, err = utils.FindFileWithSuffix(path, ".symbols.zip")

			if err != nil {
				return err
			}

			if aabPath == "" {
				aabPath, err = utils.FindFileWithSuffix(path, ".aab")

				if err != nil {
					return err
				}
			}
		} else if strings.HasSuffix(path, ".symbols.zip") {
			zipPath = path

			if aabPath == "" {
				buildDirectory := filepath.Dir(path)

				aabPath, err = utils.FindFileWithSuffix(buildDirectory, ".aab")

				if err != nil {
					return err
				}
			}
		} else {
			return fmt.Errorf("unsupported files. Please specify the unity `symbols.zip` file.")
		}
	}

	log.Info("Using " + aabPath + " as the Unity Android AAB file")

	log.Info("Using " + zipPath + " as the Unity Android symbols zip file")

	tempAabDir, err := os.MkdirTemp("", "bugsnag-cli-unity-android-aab-unpacking-*")

	if err != nil {
		return fmt.Errorf("error creating temporary working directory " + err.Error())
	}

	defer os.RemoveAll(tempAabDir)

	log.Info("Extracting " + filepath.Base(aabPath) + " to " + tempAabDir)

	err = utils.Unzip(aabPath, tempAabDir)

	if err != nil {
		return err
	}

	aabManifestPathExpected := filepath.Join(tempAabDir, "base", "manifest", "AndroidManifest.xml")
	if utils.FileExists(aabManifestPathExpected) {
		aabManifestPath = aabManifestPathExpected
	} else {
		log.Warn("AndroidManifest.xml not found in AAB file")
	}

	if aabManifestPath != "" && (applicationId == "" || buildUuid == "" || versionCode == "" || versionName == "") {

		log.Info("Reading data from AndroidManifest.xml")

		manifestData, err = android.ReadAabManifest(filepath.Join(aabManifestPath))

		if err != nil {
			return fmt.Errorf("unable to read data from " + aabManifestPath + " " + err.Error())
		}

		if applicationId == "" {
			applicationId = manifestData["applicationId"]
			if applicationId != "" {
				log.Info("Using " + applicationId + " as application ID from AndroidManifest.xml")
			}
		}

		if buildUuid == "" {
			buildUuid = manifestData["buildUuid"]
			if buildUuid != "" {
				log.Info("Using " + buildUuid + " as build ID from AndroidManifest.xml")
			} else {
				buildUuid = android.GetDexBuildId(filepath.Join(tempAabDir, "base", "dex"))

				if buildUuid != "" {
					log.Info("Using " + buildUuid + " as build ID from dex signatures")
				}
			}
		} else if buildUuid == "none" {
			log.Info("No build ID will be used")
			buildUuid = ""
		}

		if versionCode == "" {
			versionCode = manifestData["versionCode"]
			if versionCode != "" {
				log.Info("Using " + versionCode + " as version code from AndroidManifest.xml")
			}
		}

		if versionName == "" {
			versionName = manifestData["versionName"]
			if versionName != "" {
				log.Info("Using " + versionName + " as version name from AndroidManifest.xml")
			}
		}
	}

	soFilePath := filepath.Join(tempAabDir, "BUNDLE-METADATA", "com.android.tools.build.debugsymbols")

	fileList, err := utils.BuildFileList([]string{soFilePath})

	if len(fileList) > 0 && err == nil {
		for _, file := range fileList {
			err = ProcessAndroidNDK(apiKey, applicationId, "", "", []string{file}, projectRoot, "", versionCode, versionName, endpoint, failOnUploadError, retries, timeout, overwrite, dryRun)

			if err != nil {
				return err
			}
		}
	} else {
		log.Info("No NDK (.so) files detected for upload. " + err.Error())
	}

	mappingFilePath := filepath.Join(tempAabDir, "BUNDLE-METADATA", "com.android.tools.build.obfuscation", "proguard.map")

	if utils.FileExists(mappingFilePath) {
		err = ProcessAndroidProguard(apiKey, applicationId, "", buildUuid, []string{mappingFilePath}, "", versionCode, versionName, endpoint, retries, timeout, overwrite, dryRun)

		if err != nil {
			return err
		}
	} else {
		log.Info("No Proguard (mapping.txt) file detected for upload.")
	}

	log.Info("Processing Unity Android symbol files")

	tempDir, err := os.MkdirTemp("", "bugsnag-cli-unity-android-*")

	if err != nil {
		return fmt.Errorf("error creating temporary working directory " + err.Error())
	}

	defer os.RemoveAll(tempDir)

	log.Info("Extracting " + filepath.Base(zipPath) + " to " + tempDir)

	err = utils.Unzip(zipPath, tempDir)

	if err != nil {
		return err
	}

	if arch == "" {
		archList, err = utils.BuildFolderList([]string{tempDir})
		if err != nil {
			return err
		}
	} else {
		archList = []string{arch}
	}

	for _, arch := range archList {
		soPath := filepath.Join(tempDir, arch)
		fileList, err := utils.BuildFileList([]string{soPath})
		if err != nil {
			return err
		}
		for _, file := range fileList {
			if filepath.Base(file) == "libil2cpp.sym.so" && utils.ContainsString(fileList, "libil2cpp.dbg.so") {
				continue
			}
			symbolFileList = append(symbolFileList, file)
		}
	}

	numberOfFiles := len(symbolFileList)

	if numberOfFiles < 1 {
		log.Info("No symbol files found to process")
		return nil
	}

	for _, file := range symbolFileList {
		uploadOptions, err := utils.BuildAndroidNDKUploadOptions(apiKey, applicationId, versionName, versionCode, projectRoot, filepath.Base(file), overwrite)

		if err != nil {
			return err
		}

		fileFieldData := make(map[string]string)
		fileFieldData["soFile"] = file

		err = server.ProcessRequest(endpoint+"/ndk-symbol", uploadOptions, fileFieldData, timeout, file, dryRun)

		if err != nil {
			if numberOfFiles > 1 && failOnUploadError {
				return err
			} else {
				log.Warn(err.Error())
			}
		} else {
			log.Success("Uploaded " + filepath.Base(file))
		}
	}

	return nil
}
