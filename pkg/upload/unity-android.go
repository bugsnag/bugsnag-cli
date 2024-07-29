package upload

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bugsnag/bugsnag-cli/pkg/android"
	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
)

type UnityAndroid struct {
	Path          utils.Paths `arg:"" name:"path" help:"The path to the Unity symbols (.zip) file to upload (or directory containing it)" type:"path"`
	AabPath       utils.Path  `help:"The path to an AAB file to upload alongside the Unity symbols"`
	ApplicationId string      `help:"A unique application ID, usually the package name, of the application"`
	BuildUuid     string      `help:"A unique identifier for this build of the application" xor:"no-build-uuid,build-uuid"`
	NoBuildUuid   bool        `help:"Prevents the automatically generated build UUID being uploaded with the build" xor:"build-uuid,no-build-uuid"`
	ProjectRoot   string      `help:"The path to strip from the beginning of source file names referenced in stacktraces on the BugSnag dashboard" type:"path"`
	VersionCode   string      `help:"The version code of this build of the application"`
	VersionName   string      `help:"The version of the application"`
}

func ProcessUnityAndroid(
	apiKey string,
	options UnityAndroid,
	endpoint string,
	timeout int,
	retries int,
	overwrite bool,
	dryRun bool,
	logger log.Logger,
) error {
	var err error
	var zipPath string
	var archList []string
	var symbolFileList []string
	var manifestData map[string]string
	var aabPath = string(options.AabPath)

	for _, path := range options.Path {
		if utils.IsDir(path) {
			zipPath, err = utils.FindLatestFileWithSuffix(path, ".symbols.zip")

			if err != nil {
				return err
			}

			if aabPath == "" {
				aabPath, _ = utils.FindLatestFileWithSuffix(path, ".aab")
			}
		} else if strings.HasSuffix(path, ".symbols.zip") {
			zipPath = path

			if aabPath == "" {
				buildDirectory := filepath.Dir(path)

				aabPath, _ = utils.FindLatestFileWithSuffix(buildDirectory, ".aab")
			}
		} else {
			return fmt.Errorf("%s is not a .symbols.zip file or containing directory", path)
		}
	}

	if aabPath != "" {

		logger.Debug(fmt.Sprintf("Extracting %s into a temporary directory", filepath.Base(aabPath)))

		aabDir, err := utils.ExtractFile(aabPath, "aab")

		if err != nil {
			return err
		}

		defer os.RemoveAll(aabDir)

		manifestData, err = android.MergeUploadOptionsFromAabManifest(aabDir, apiKey, options.ApplicationId, options.BuildUuid, options.NoBuildUuid, options.VersionCode, options.VersionName, logger)

		if err != nil {
			return err
		}

		err = ProcessAndroidAab(
			manifestData["apiKey"],
			AndroidAabMapping{
				ApplicationId: manifestData["applicationId"],
				BuildUuid:     manifestData["buildUuid"],
				NoBuildUuid:   options.NoBuildUuid,
				Path:          []string{aabDir},
				ProjectRoot:   options.ProjectRoot,
				VersionCode:   manifestData["versionCode"],
				VersionName:   manifestData["versionName"],
			},
			endpoint,
			retries,
			timeout,
			overwrite,
			dryRun,
			logger,
		)

		if err != nil {
			return err
		}
	}

	logger.Debug(fmt.Sprintf("Extracting %s into a temporary directory", filepath.Base(zipPath)))

	if manifestData == nil {
		manifestData, _ = android.MergeUploadOptionsFromAabManifest("", apiKey, options.ApplicationId, options.BuildUuid, options.NoBuildUuid, options.VersionCode, options.VersionName, logger)
	}

	unityDir, err := utils.ExtractFile(zipPath, "unity-android")

	if err != nil {
		return err
	}

	defer os.RemoveAll(unityDir)

	archList, err = utils.BuildDirectoryList([]string{unityDir})

	if err != nil {
		return err
	}

	for _, arch := range archList {
		soPath := filepath.Join(unityDir, arch)
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

	err = android.UploadAndroidNdk(
		symbolFileList,
		manifestData["apiKey"],
		manifestData["applicationId"],
		manifestData["versionName"],
		manifestData["versionCode"],
		options.ProjectRoot,
		overwrite,
		endpoint,
		timeout,
		retries,
		dryRun,
		logger,
	)

	if err != nil {
		return err
	}

	return nil
}
