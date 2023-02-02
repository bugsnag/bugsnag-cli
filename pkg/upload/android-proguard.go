package upload

import (
	"path/filepath"
	"strconv"

	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/server"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
)

type AndroidProguardMapping struct {
	ApplicationId   string            `help:"Module application identifier"`
	AppManifestPath string            `help:"(required) Path to directory or file to upload" type:"path"`
	BuildUuid       string            `help:"Module Build UUID"`
	Configuration   string            `help:"Build type, like 'debug' or 'release'"`
	Path            utils.UploadPaths `arg:"" name:"path" help:"(required) Path to directory or file to upload" type:"path"`
	VersionCode     string            `help:"Module version code"`
	VersionName     string            `help:"Module version name"`
	DryRun          bool              `help:"Validate but do not upload"`
}

func ProcessAndroidProguard(paths []string, applicationId string, appManifestPath string, buildUuid string, configuration string, versionCode string, versionName string, endpoint string, timeout int, retries int, overwrite bool, apiKey string, failOnUploadError bool, dryRun bool) error {

	uploadFileOptions := make(map[string]map[string]string)

	for _, path := range paths {
		// Path is a directory
		if utils.IsDir(path) {
			mergedManifestsPath := filepath.Join(path, "build", "intermediates", "merged_manifests")
			if utils.IsDir(mergedManifestsPath) {
				variants, err := utils.BuildVariantsList(mergedManifestsPath)
				if err != nil {
					log.Error(err.Error(), 1)
				}

				for _, variant := range variants {
					uploadFileOptions[variant] = map[string]string{}
					uploadFileOptions[variant]["androidManifestPath"] = filepath.Join(mergedManifestsPath, variant, "AndroidManifest.xml")
					uploadFileOptions[variant]["mappingPath"] = filepath.Join(mergedManifestsPath, "..", "..", "outputs", "mapping", variant, "mapping.txt")
				}
			} else {
				log.Error("unable to find `merged_manifests` in "+path, 1)
			}
			//	Path is a file
		} else {

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

			if filepath.Base(path) == "mapping.txt" {
				uploadFileOptions[configuration] = map[string]string{}
				uploadFileOptions[configuration]["androidManifestPath"] = appManifestPath
				uploadFileOptions[configuration]["mappingPath"] = path
			} else {
				log.Error(path+" is not a supported file. Please use `mapping.txt`", 1)
			}
		}
	}

	numberOfVariants := len(uploadFileOptions)

	if numberOfVariants < 1 {
		log.Info("No variants to process")
		return nil
	}

	log.Info("Processing " + strconv.Itoa(numberOfVariants) + " variant(s)")

	for variant, config := range uploadFileOptions {
		log.Info("Processing mapping.txt for variant: " + variant)

		log.Info("Gathering information from AndroidManifest.xml")

		androidManifestData, err := utils.ParseAndroidManifestXML(config["androidManifestPath"])

		if err != nil {
			return err
		}

		if buildUuid == "" {
			for i := range androidManifestData.Application.MetaData.Name {
				if androidManifestData.Application.MetaData.Name[i] == "com.bugsnag.android.BUILD_UUID" {
					buildUuid = androidManifestData.Application.MetaData.Value[i]
				}
			}
		}

		log.Info("Compressing " + config["mappingPath"])

		outputFile, err := utils.GzipCompress(config["mappingPath"])

		if err != nil {
			return err
		}

		log.Info("Uploading debug information for " + config["mappingPath"])

		uploadOptions := utils.BuildAndroidProguardUploadOptions(apiKey, androidManifestData.AppId, androidManifestData.VersionName, androidManifestData.VersionCode, buildUuid, overwrite)

		requestStatus := server.ProcessRequest(endpoint, uploadOptions, "proguard", outputFile, timeout)

		if requestStatus != nil {
			if numberOfVariants > 1 && failOnUploadError {
				return requestStatus
			} else {
				log.Warn(requestStatus.Error())
			}
		} else {
			log.Success(config["mappingPath"] + " uploaded")
		}

	}

	return nil
}
