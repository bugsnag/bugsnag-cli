package upload

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/bugsnag/bugsnag-cli/pkg/android"
	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
)

type AndroidAabMapping struct {
	ApplicationId string            `help:"Module application identifier"`
	BuildUuid     string            `help:"Module Build UUID"`
	Configuration string            `help:"Build type, like 'debug' or 'release'"`
	Path          utils.UploadPaths `arg:"" name:"path" help:"(required) Path to directory or file to upload" type:"path"`
	ProjectRoot   string            `help:"path to remove from the beginning of the filenames in the mapping file"`
	VersionCode   string            `help:"Module version code"`
	VersionName   string            `help:"Module version name"`
}

var getApiKeyFromManifest = false
var variantConfig = make(map[string]map[string]string)

func ProcessAndroidAab(paths []string, buildUuid string, configuration string, projectRoot string, versionCode string, versionName string, endpoint string, timeout int, retries int, overwrite bool, apiKey string, failOnUploadError bool) error {
	// Check if we have project root
	if projectRoot == "" {
		return fmt.Errorf("`--project-root` missing from options")
	}

	// Check if we need to get the API key from the AndroidManifest.xml
	if apiKey == "" {
		getApiKeyFromManifest = true
	}

	// Loop through paths
	for _, path := range paths {
		log.Info("Building variants configuration")

		// Check to see if we're working with a directory
		if utils.IsDir(path) {
			// Build a list of variants from the directory
			variants, err := android.BuildVariantsList(path)

			if err != nil {
				return err
			}

			// For each variant build a map that has the path to the .aab file
			for _, variant := range variants {
				variantConfig[variant] = map[string]string{}
				variantConfig[variant]["aabPath"] = filepath.Join(path, variant, "app-"+variant+".aab")
			}

			// Check to see if we're working with a single AAB file
		} else if filepath.Ext(path) == ".aab" {

			variant := filepath.Base(strings.Replace(path, filepath.Base(path), "", -1))

			variantConfig[variant] = map[string]string{}

			variantConfig[variant]["aabPath"] = path

		} else {
			return fmt.Errorf("unsupported file type for " + path)
		}
	}

	numberOfVariants := len(variantConfig)

	if numberOfVariants < 1 {
		log.Info("No variants to process")
		return nil
	}

	log.Info("Processing " + strconv.Itoa(numberOfVariants) + " variant(s)")

	for variant, config := range variantConfig {
		log.Info("Processing variant: " + variant)

		if !utils.FileExists(config["aabPath"]) {
			log.Info(variant + " does not have a valid .aab file")
			continue
		}

		outputPath := filepath.Join(strings.Replace(config["aabPath"], filepath.Base(config["aabPath"]), "", -1), "raw")

		log.Info("Expanding " + filepath.Base(config["aabPath"]))

		err := utils.Unzip(config["aabPath"], outputPath)

		if err != nil {
			return err
		}

		log.Success(filepath.Base(config["aabPath"]) + " expanded")

		log.Info("Reading data from AndroidManifest.xml")

		aabManifestData, err := android.ReadAabManifest(filepath.Join(outputPath, "base", "manifest", "AndroidManifest.xml"))

		if err != nil {
			return fmt.Errorf("error reading raw AAB manifest data. " + err.Error())
		}

		if getApiKeyFromManifest {
			apiKey = aabManifestData["apiKey"]
		}

		err = android.ProcessProguard(apiKey, variant, outputPath, aabManifestData, overwrite, timeout, numberOfVariants, endpoint, failOnUploadError)

		if err != nil {
			return fmt.Errorf("error processing Proguard mapping file. " + err.Error())
		}

		androidNdkRoot, err := android.GetAndroidNDKRoot("")

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

		err = android.ProcessNdk(apiKey, variant, outputPath, aabManifestData, objCopyPath, projectRoot, overwrite, timeout, endpoint, failOnUploadError)

		if err != nil {
			return fmt.Errorf("error processing NDK symbol files. " + err.Error())
		}
	}

	return nil
}
