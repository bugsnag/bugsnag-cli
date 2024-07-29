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
	options AndroidNdkMapping,
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

	for _, path := range options.Path {

		// Search for NDK symbol files based on an expected path
		arr := []string{"android", "app", "build", "intermediates", "merged_native_libs"}
		mergeNativeLibPath, err = android.FindNativeLibPath(arr, path)

		if err != nil {
			return err
		}

		if filepath.Base(mergeNativeLibPath) == "merged_native_libs" {
			if options.Variant == "" {
				options.Variant, err = android.GetVariantDirectory(mergeNativeLibPath)
				if err != nil {
					return err
				}
			}

			if options.AppManifest == "" {
				appManifestPathExpected = filepath.Join(mergeNativeLibPath, "..", "merged_manifests", options.Variant, "AndroidManifest.xml")
				if utils.FileExists(appManifestPathExpected) {
					options.AppManifest = appManifestPathExpected
					logger.Debug(fmt.Sprintf("Found app manifest at: %s", options.AppManifest))
				}

			}

			if options.ProjectRoot == "" {
				// Setting options.ProjectRoot to the suspected root of the project
				options.ProjectRoot = filepath.Join(mergeNativeLibPath, "..", "..", "..", "..")
			}
		}

		// Ensure only files from within the directory the upload command is run from are uploaded
		if !utils.IsDir(path) {
			fileList = append(fileList, path)
		} else if strings.Contains(path, fmt.Sprintf("merged_native_libs/%s", options.Variant)) {
			fileList, err = utils.BuildFileList([]string{path})
		} else {
			fileList, err = utils.BuildFileList([]string{filepath.Join(mergeNativeLibPath, options.Variant)})
		}

		if err != nil {
			return fmt.Errorf("error building file list for options.Variant: " + options.Variant + ". " + err.Error())
		}
	}

	if options.ProjectRoot != "" {
		logger.Debug(fmt.Sprintf("Using %s as the project root", options.ProjectRoot))
	}

	// Check to see if we need to read the manifest file due to missing options
	if options.AppManifest != "" && (apiKey == "" || options.ApplicationId == "" || options.VersionCode == "" || options.VersionName == "") {

		logger.Debug("Reading data from AndroidManifest.xml")
		manifestData, err := android.ParseAndroidManifestXML(options.AppManifest)

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

		if options.ApplicationId == "" {
			options.ApplicationId = manifestData.ApplicationId

			if options.ApplicationId != "" {
				logger.Debug(fmt.Sprintf("Using %s as application ID from AndroidManifest.xml", options.ApplicationId))
			}
		}

		if options.VersionCode == "" {
			options.VersionCode = manifestData.VersionCode

			if options.VersionCode != "" {
				logger.Debug(fmt.Sprintf("Using %s as version code from AndroidManifest.xml", options.VersionCode))
			}
		}

		if options.VersionName == "" {
			options.VersionName = manifestData.VersionName

			if options.VersionName != "" {
				logger.Debug(fmt.Sprintf("Using %s as version name from AndroidManifest.xml", options.VersionName))
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
				options.AndroidNdkRoot, err = android.GetAndroidNDKRoot(options.AndroidNdkRoot)

				if err != nil {
					return err
				}

				objCopyPath, err = android.BuildObjcopyPath(options.AndroidNdkRoot)

				if err != nil {
					return err
				}
				logger.Debug(fmt.Sprintf("Located objcopy within Android NDK path: %s", options.AndroidNdkRoot))
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
		options.ApplicationId,
		options.VersionName,
		options.VersionCode,
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
