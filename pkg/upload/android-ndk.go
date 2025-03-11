package upload

import (
	"fmt"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/bugsnag/bugsnag-cli/pkg/android"
	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/options"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
)

func ProcessAndroidNDK(options options.CLI, endpoint string, logger log.Logger) error {
	ndkOptions := options.Upload.AndroidNdk

	var fileList []string
	var symbolFileList []string
	var mergeNativeLibPath string
	var err error
	var workingDir string
	var appManifestPathExpected string
	var objCopyPath string

	soFilePattern := `\.so.*$`
	soFileRegex := regexp.MustCompile(soFilePattern)

	for _, path := range ndkOptions.Path {

		// Search for NDK symbol files based on an expected path
		arr := []string{"android", "app", "build", "intermediates", "merged_native_libs"}
		mergeNativeLibPath, err = android.FindNativeLibPath(arr, path)

		if err != nil {
			return err
		}

		if filepath.Base(mergeNativeLibPath) == "merged_native_libs" {
			if ndkOptions.Variant == "" {
				ndkOptions.Variant, err = android.GetVariantDirectory(mergeNativeLibPath)
				if err != nil {
					return err
				}
			}

			if ndkOptions.AppManifest == "" {
				logger.Info("No app manifest provided, attempting to find one")
				appManifestPathExpected = filepath.Join(path, "app", "build", "intermediates", "merged_manifests", ndkOptions.Variant, "AndroidManifest.xml")
				logger.Info(fmt.Sprintf("Looking for app manifest at: %s", appManifestPathExpected))
				if utils.FileExists(appManifestPathExpected) {
					ndkOptions.AppManifest = appManifestPathExpected
					logger.Debug(fmt.Sprintf("Found app manifest at: %s", ndkOptions.AppManifest))
				} else {
					appManifestPathExpected = filepath.Join(path, "app", "build", "intermediates", "merged_manifests", ndkOptions.Variant, "process"+cases.Title(language.English).String(ndkOptions.Variant)+"Manifest", "AndroidManifest.xml")
					logger.Info(fmt.Sprintf("Looking for app manifest at: %s", appManifestPathExpected))
					if utils.FileExists(appManifestPathExpected) {
						ndkOptions.AppManifest = appManifestPathExpected
						logger.Info(fmt.Sprintf("Found app manifest at: %s", ndkOptions.AppManifest))
					} else {
						logger.Info(fmt.Sprintf("No app manifest found at: %s", appManifestPathExpected))
					}
				}
			}

			if ndkOptions.ProjectRoot == "" {
				// Setting options.ProjectRoot to the suspected root of the project
				ndkOptions.ProjectRoot = filepath.Join(mergeNativeLibPath, "..", "..", "..", "..")
			}
		}

		// Ensure only files from within the directory the upload command is run from are uploaded
		if !utils.IsDir(path) {
			fileList = append(fileList, path)
		} else if strings.Contains(path, fmt.Sprintf("merged_native_libs/%s", ndkOptions.Variant)) {
			fileList, err = utils.BuildFileList([]string{path})
		} else {
			fileList, err = utils.BuildFileList([]string{filepath.Join(mergeNativeLibPath, ndkOptions.Variant)})
		}

		if err != nil {
			return fmt.Errorf("error building file list for options.Variant: %s. %w", ndkOptions.Variant, err)
		}
	}

	if ndkOptions.ProjectRoot != "" {
		logger.Debug(fmt.Sprintf("Using %s as the project root", ndkOptions.ProjectRoot))
	}

	// Check to see if we need to read the manifest file due to missing options
	if ndkOptions.AppManifest != "" && (options.ApiKey == "" || ndkOptions.ApplicationId == "" || ndkOptions.VersionCode == "" || ndkOptions.VersionName == "") {

		logger.Debug("Reading data from AndroidManifest.xml")
		manifestData, err := android.ParseAndroidManifestXML(ndkOptions.AppManifest)

		if err != nil {
			return err
		}

		if options.ApiKey == "" {
			for key, value := range manifestData.Application.MetaData.Name {
				if value == "com.bugsnag.android.API_KEY" {
					options.ApiKey = manifestData.Application.MetaData.Value[key]
				}
			}

			if options.ApiKey != "" {
				logger.Debug(fmt.Sprintf("Using %s as API key from AndroidManifest.xml", options.ApiKey))
			}
		}

		if ndkOptions.ApplicationId == "" {
			ndkOptions.ApplicationId = manifestData.ApplicationId

			if ndkOptions.ApplicationId != "" {
				logger.Debug(fmt.Sprintf("Using %s as application ID from AndroidManifest.xml", ndkOptions.ApplicationId))
			}
		}

		if ndkOptions.VersionCode == "" {
			ndkOptions.VersionCode = manifestData.VersionCode

			if ndkOptions.VersionCode != "" {
				logger.Debug(fmt.Sprintf("Using %s as version code from AndroidManifest.xml", ndkOptions.VersionCode))
			}
		}

		if ndkOptions.VersionName == "" {
			ndkOptions.VersionName = manifestData.VersionName

			if ndkOptions.VersionName != "" {
				logger.Debug(fmt.Sprintf("Using %s as version name from AndroidManifest.xml", ndkOptions.VersionName))
			}
		}
	}

	// Process .so files through objcopy to create .sym files, filtering any other file type
	for _, file := range fileList {
		if strings.HasSuffix(file, ".so.sym") {
			symbolFileList = append(symbolFileList, file)
		} else if soFileRegex.MatchString(file) {
			// Check NDK path is set
			if objCopyPath == "" {
				ndkOptions.AndroidNdkRoot, err = android.GetAndroidNDKRoot(ndkOptions.AndroidNdkRoot)

				if err != nil {
					return err
				}

				objCopyPath, err = android.BuildObjcopyPath(ndkOptions.AndroidNdkRoot)

				if err != nil {
					return err
				}
				logger.Debug(fmt.Sprintf("Located objcopy within Android NDK path: %s", ndkOptions.AndroidNdkRoot))
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
		options.ApiKey,
		ndkOptions.ApplicationId,
		ndkOptions.VersionName,
		ndkOptions.VersionCode,
		ndkOptions.ProjectRoot,
		endpoint,
		options,
		logger,
	)

	if err != nil {
		return err
	}

	return nil
}
