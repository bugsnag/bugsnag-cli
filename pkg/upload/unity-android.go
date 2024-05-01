package upload

import (
	"fmt"
	"github.com/bugsnag/bugsnag-cli/pkg/android"
	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
	"os"
	"path/filepath"
	"strings"
)

type UnityAndroid struct {
	AabPath       utils.Path  `help:"Path to Android AAB file to upload with your Unity symbols"`
	ApplicationId string      `help:"Module application identifier"`
	Path          utils.Paths `arg:"" name:"path" help:"(required) Path to Unity symbols zip file or directory to upload" type:"path"`
	ProjectRoot   string      `help:"path to remove from the beginning of the filenames in the mapping file" type:"path"`
	VersionCode   string      `help:"Module version code"`
	VersionName   string      `help:"Module version name"`
	BuildUuid     string      `help:"Module Build UUID" xor:"no-build-uuid,build-uuid"`
	NoBuildUuid   bool        `help:"Upload with no Build UUID" xor:"build-uuid,no-build-uuid"`
}

func ProcessUnityAndroid(
	apiKey string,
	aabPath string,
	applicationId string,
	versionCode string,
	buildUuid string,
	noBuildUuid bool,
	versionName string,
	projectRoot string,
	paths []string,
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

	for _, path := range paths {
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

		logger.Info(fmt.Sprintf("Extracting %s into a temporary directory", filepath.Base(aabPath)))

		aabDir, err := utils.ExtractFile(aabPath, "aab")

		if err != nil {
			return err
		}

		defer os.RemoveAll(aabDir)

		manifestData, err = android.MergeUploadOptionsFromAabManifest(aabDir, apiKey, applicationId, buildUuid, noBuildUuid, versionCode, versionName, logger)

		if err != nil {
			return err
		}

		err = ProcessAndroidAab(
			manifestData["apiKey"],
			manifestData["applicationId"],
			manifestData["buildUuid"],
			noBuildUuid,
			[]string{aabDir},
			projectRoot,
			manifestData["versionCode"],
			manifestData["versionName"],
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

	logger.Info(fmt.Sprintf("Extracting %s into a temporary directory", filepath.Base(zipPath)))

	if manifestData == nil {
		manifestData, _ = android.MergeUploadOptionsFromAabManifest("", apiKey, applicationId, buildUuid, noBuildUuid, versionCode, versionName, logger)
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
		projectRoot,
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
