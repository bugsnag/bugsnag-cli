package upload

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/bugsnag/bugsnag-cli/pkg/android"
	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
)

type AndroidAabMapping struct {
	AndroidNdkRoot string            `help:"Path to Android NDK installation ($ANDROID_NDK_ROOT)"`
	ApplicationId  string            `help:"Module application identifier"`
	BuildUuid      string            `help:"Module Build UUID"`
	Configuration  string            `help:"Build type, like 'debug' or 'release'"`
	DryRun         bool              `help:"Validate but do not upload"`
	MappingPath    string            `help:"Path to app mapping file"`
	Path           utils.UploadPaths `arg:"" name:"path" help:"(required) Path to directory or file to upload" type:"path"`
	ProjectRoot    string            `help:"path to remove from the beginning of the filenames in the mapping file" type:"path"`
	VersionCode    string            `help:"Module version code"`
	VersionName    string            `help:"Module version name"`
}

func ProcessAndroidAab(appId string, buildUuid string, configuration string, paths []string, projectRoot string, versionCode string, versionName string, androidNdkRoot string, mappingPath string, dryRun bool, endpoint string, timeout int, retries int, overwrite bool, apiKey string, failOnUploadError bool) error {

	var getApiKeyFromManifest = false
	var aabPath string

	for _, path := range paths {

		if projectRoot == "" {
			projectRoot = path
		}

		if apiKey == "" {
			getApiKeyFromManifest = true
		}

		if utils.FileExists(filepath.Join(path, "app", "build", "outputs", "bundle")) {
			if configuration == "" {
				variants, err := android.BuildVariantsList(filepath.Join(path, "app", "build", "outputs", "bundle"))

				if err != nil {
					return err
				}

				if len(variants) > 1 {
					fmt.Println(variants)
					return fmt.Errorf("more than one variant found. Please specify using `--configuration`")
				}

				configuration = variants[0]
			}
			aabPath = filepath.Join(path, "app", "build", "outputs", "bundle", configuration, "app-"+configuration+".aab")

			if mappingPath == "" {
				if utils.FileExists(filepath.Join(path, "app", "build", "outputs", "mapping", configuration, "mapping.txt")) {
					mappingPath = filepath.Join(path, "app", "build", "outputs", "mapping", configuration, "mapping.txt")
				} else {
					return fmt.Errorf("unable to find mapping.txt. Please specify using `--mapping-path` ")
				}
			}

		} else if filepath.Ext(path) == ".aab" {
			if configuration == "" {
				return fmt.Errorf("missing configuration. Please specify using `--configuration`")
			} else if mappingPath == "" {
				return fmt.Errorf("unable to find mapping.txt. Please specify using `--mapping-path` ")
			} else {
				aabPath = path
			}
		} else {
			return fmt.Errorf("")
		}

		log.Info("Processing app-" + configuration + ".aab")

		if !utils.FileExists(aabPath) {
			return fmt.Errorf(aabPath + " does not exist on the system.")
		}

		outputPath := filepath.Join(strings.Replace(aabPath, filepath.Base(aabPath), "", -1), "raw")

		log.Info("Expanding " + filepath.Base(aabPath))

		err := utils.Unzip(aabPath, outputPath)

		if err != nil {
			return err
		}

		log.Success(filepath.Base(aabPath) + " expanded")

		aabManifestData, err := android.ReadAabManifest(filepath.Join(outputPath, "base", "manifest", "AndroidManifest.xml"))

		if err != nil {
			return fmt.Errorf("error reading raw AAB manifest data. " + err.Error())
		}

		if appId == "" {
			appId = aabManifestData["package"]
		}

		if buildUuid == "" {
			buildUuid = aabManifestData["buildUuid"]
		}

		if versionCode == "" {
			versionCode = aabManifestData["versionCode"]
		}

		if versionName == "" {
			versionName = aabManifestData["versionName"]
		}

		if getApiKeyFromManifest {
			apiKey = aabManifestData["apiKey"]
		}

		err = android.ProcessProguard(apiKey, configuration, outputPath, appId, versionCode, versionName, buildUuid, mappingPath, overwrite, timeout, endpoint, failOnUploadError)

		if err != nil {
			return fmt.Errorf("error processing Proguard mapping file. " + err.Error())
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

		err = android.ProcessNdk(apiKey, configuration, outputPath, appId, versionCode, versionName, objCopyPath, projectRoot, overwrite, timeout, endpoint, failOnUploadError)

		if err != nil {
			return fmt.Errorf("error processing NDK symbol files. " + err.Error())
		}
	}

	return nil
}
