package upload

import (
	"fmt"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/bugsnag/bugsnag-cli/pkg/log"
	//"github.com/bugsnag/bugsnag-cli/pkg/server"
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

var packageName string

func ProcessAndroidAab(paths []string, buildUuid string, configuration string, projectRoot string, versionCode string, versionName string, endpoint string, timeout int, retries int, overwrite bool, apiKey string, failOnUploadError bool) error {
	// Check if we have project root
	if projectRoot == "" {
		return fmt.Errorf("`--project-root` missing from options")
	}

	var setApiKey = false

	if apiKey == "" {
		setApiKey = true
	}

	variantConfig := make(map[string]map[string]string)

	for _, path := range paths {
		if utils.IsDir(path) {
			// build a list of variants
			variants, err := utils.BuildVariantsList(path)
			for _, variant := range variants {
				variantConfig[variant] = map[string]string{}
				variantConfig[variant]["aabPath"] = filepath.Join(path, variant, "app-"+variant+".aab")
			}

			if err != nil {
				return err
			}

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

	for variant, config := range variantConfig {
		log.Info("Unzipping aab file for " + variant)
		if !utils.FileExists(config["aabPath"]) {
			log.Info(variant + " does not have a valid .aab file")
			continue
		}
		outputPath := filepath.Join(strings.Replace(config["aabPath"], filepath.Base(config["aabPath"]), "", -1), "raw")
		err := utils.Unzip(config["aabPath"], outputPath)
		if err != nil {
			return err
		}
		log.Success(config["aabPath"] + " unzipped")

		log.Info("Reading data from AndroidManifest.xml")

		aabManifestData, err := utils.ReadAabManifest(filepath.Join(outputPath, "base", "manifest", "AndroidManifest.xml"))

		if err != nil {
			return fmt.Errorf("error reading raw AAB manifest data. " + err.Error())
		}

		if setApiKey {
			apiKey = aabManifestData["apiKey"]
		}

		err = AabProcessProguard(variant,outputPath, aabManifestData, overwrite,timeout,numberOfVariants)

		if err != nil {
			return fmt.Errorf("error processing Proguard mapping file. " + err.Error())
		}


		//symbolPath := []string{filepath.Join(outputPath,"BUNDLE-METADATA", "com.android.tools.build.debugsymbols")}

	}

	return nil
}

func AabProcessProguard(variant string, outputPath string, aabManifestData map[string]string, overwrite bool, timeout int, numberOfVariants int)error{
	log.Info("Processing Proguard mapping for " + variant)

	proguardMappingPath := filepath.Join(outputPath, "BUNDLE-METADATA", "com.android.tools.build.obfuscation", "proguard.map")

	if !utils.FileExists(proguardMappingPath) {
		return fmt.Errorf(proguardMappingPath + " does not exist")
	}

	log.Info("Compressing " + proguardMappingPath)

	outputFile, err := utils.GzipCompress(proguardMappingPath)

	if err != nil {
		return err
	}

	log.Info("Uploading debug information for " + outputFile)

	uploadOptions := utils.BuildAndroidProguardUploadOptions(apiKey, aabManifestData["package"], aabManifestData["versionName"], aabManifestData["versionCode"], aabManifestData["buildUuid"], overwrite)

	requestStatus := server.ProcessRequest(endpoint, uploadOptions, "proguard", outputFile, timeout)

	if requestStatus != nil {
		if numberOfVariants > 1 && failOnUploadError {
			return requestStatus
		} else {
			log.Warn(requestStatus.Error())
		}
	} else {
		log.Success(proguardMappingPath + " uploaded")
	}

	return nil
}
