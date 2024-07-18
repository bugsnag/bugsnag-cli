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

type AndroidNdkMapping struct {
	Path           utils.Paths `arg:"" name:"path" help:"The path to the directory or file to upload" type:"path" default:"."`
	ApplicationId  string      `help:"A unique application ID, usually the package name, of the application"`
	AndroidNdkRoot string      `help:"The path to your NDK installation, used to access the objcopy tool for extracting symbol information"`
	AppManifest    string      `help:"The path to a manifest file (AndroidManifest.xml) from which to obtain build information" type:"path"`
	ProjectRoot    string      `help:"The path to strip from the beginning of source file names referenced in stacktraces on the BugSnag dashboard" type:"path"`
	Variant        string      `help:"The build type/flavor (e.g. debug, release) used to disambiguate the between built files when searching the project directory"`
	VersionCode    string      `help:"The version code of this build of the application"`
	VersionName    string      `help:"The version of the application"`
}

func ProcessAndroidNDK(
	apiKey string,
	applicationId string,
	androidNdkRoot string,
	appManifestPath string,
	paths []string,
	projectRoot string,
	variant string,
	versionCode string,
	versionName string,
	endpoint string,
	retries int,
	timeout int,
	overwrite bool,
	dryRun bool,
	logger log.Logger,
) error {

	var fileList []string
	var symbolFileList []string
	var mergeNativeLibPath string
	var err error
	var workingDir string
	var appManifestPathExpected string
	var objCopyPath string

	for _, path := range paths {

		// Search for NDK symbol files based on an expected path
		arr := []string{"android", "app", "build", "intermediates", "merged_native_libs"}
		mergeNativeLibPath, err = android.FindNativeLibPath(arr, path)

		if err != nil {
			return err
		}

		if filepath.Base(mergeNativeLibPath) == "merged_native_libs" {
			if variant == "" {
				variant, err = android.GetVariantDirectory(mergeNativeLibPath)
				if err != nil {
					return err
				}
			}

			if appManifestPath == "" {
				appManifestPathExpected = filepath.Join(mergeNativeLibPath, "..", "merged_manifests", variant, "AndroidManifest.xml")
				if utils.FileExists(appManifestPathExpected) {
					appManifestPath = appManifestPathExpected
					logger.Debug(fmt.Sprintf("Found app manifest at: %s", appManifestPath))
				}

			}

			if projectRoot == "" {
				// Setting projectRoot to the suspected root of the project
				projectRoot = filepath.Join(mergeNativeLibPath, "..", "..", "..", "..")
			}
		}

		// Ensure only files from within the directory the upload command is run from are uploaded
		if !utils.IsDir(path) {
			fileList = append(fileList, path)
		} else if strings.Contains(path, fmt.Sprintf("merged_native_libs/%s", variant)) {
			fileList, err = utils.BuildFileList([]string{path})
		} else {
			fileList, err = utils.BuildFileList([]string{filepath.Join(mergeNativeLibPath, variant)})
		}

		if err != nil {
			return fmt.Errorf("error building file list for variant: " + variant + ". " + err.Error())
		}
	}

	if projectRoot != "" {
		logger.Debug(fmt.Sprintf("Using %s as the project root", projectRoot))
	}

	// Check to see if we need to read the manifest file due to missing options
	if appManifestPath != "" && (apiKey == "" || applicationId == "" || versionCode == "" || versionName == "") {

		logger.Debug("Reading data from AndroidManifest.xml")
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

			if apiKey != "" {
				logger.Debug(fmt.Sprintf("Using %s as API key from AndroidManifest.xml", apiKey))
			}
		}

		if applicationId == "" {
			applicationId = manifestData.ApplicationId

			if applicationId != "" {
				logger.Debug(fmt.Sprintf("Using %s as application ID from AndroidManifest.xml", applicationId))
			}
		}

		if versionCode == "" {
			versionCode = manifestData.VersionCode

			if versionCode != "" {
				logger.Debug(fmt.Sprintf("Using %s as version code from AndroidManifest.xml", versionCode))
			}
		}

		if versionName == "" {
			versionName = manifestData.VersionName

			if versionName != "" {
				logger.Debug(fmt.Sprintf("Using %s as version name from AndroidManifest.xml", versionName))
			}
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
				logger.Debug(fmt.Sprintf("Located objcopy within Android NDK path: %s", androidNdkRoot))
			}

			logger.Debug(fmt.Sprintf("Extracting debug info from %s using objcopy", filepath.Base(file)))

			if workingDir == "" {
				workingDir, err = os.MkdirTemp("", "bugsnag-cli-ndk-*")

				if err != nil {
					return fmt.Errorf("error creating temporary working directory %s", err.Error())
				}

				defer os.RemoveAll(workingDir)
			}

			outputFile, err := android.Objcopy(objCopyPath, file, workingDir)

			if err != nil {
				return fmt.Errorf("failed to process file, %s using objcopy : %s", file, err.Error())
			}

			symbolFileList = append(symbolFileList, outputFile)
		}
	}

	err = android.UploadAndroidNdk(
		symbolFileList,
		apiKey,
		applicationId,
		versionName,
		versionCode,
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
