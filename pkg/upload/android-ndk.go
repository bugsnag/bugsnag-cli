package upload

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/bugsnag/bugsnag-cli/pkg/android"
	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/options"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
)

// ProcessAndroidNDK handles the processing and uploading of Android NDK (.so) symbol files.
//
// This function locates native library symbol files, extracts debug information
// using objcopy, optionally parses AndroidManifest.xml for missing metadata,
// and finally uploads the processed symbol files.
//
// Parameters:
//   - options: CLI options containing upload settings and metadata.
//   - logger: Logger instance for debug/info/error output.
//
// Returns:
//   - error: non-nil if any step in the processing or upload fails.
func ProcessAndroidNDK(options options.CLI, logger log.Logger) error {
	ndkOptions := options.Upload.AndroidNdk

	var fileList []string         // List of native library files to process
	var symbolFileList []string   // List of extracted symbol files ready for upload
	var mergeNativeLibPath string // Path to merged native libs
	var err error
	var workingDir string // Temporary working directory for objcopy output
	var appManifestPathExpected string
	var objCopyPath string // Path to objcopy executable

	// Regex to identify .so files (native libs)
	soFilePattern := `\.so.*$`
	soFileRegex := regexp.MustCompile(soFilePattern)

	for _, path := range ndkOptions.Path {
		// Locate merged native libs directory from expected Android build paths
		arr := []string{"android", "app", "build", "intermediates", "merged_native_libs"}
		mergeNativeLibPath, err = android.FindNativeLibPath(arr, path)
		if err != nil {
			return err
		}

		if filepath.Base(mergeNativeLibPath) == "merged_native_libs" {
			// Determine variant directory if not set
			if ndkOptions.Variant == "" {
				ndkOptions.Variant, err = android.GetVariantDirectory(mergeNativeLibPath)
				if err != nil {
					return err
				}
			}

			// Attempt to locate AndroidManifest.xml if not explicitly set
			if ndkOptions.AppManifest == "" {
				appManifestPathExpected = filepath.Join(mergeNativeLibPath, "..", "merged_manifests", ndkOptions.Variant, "AndroidManifest.xml")
				if utils.FileExists(appManifestPathExpected) {
					ndkOptions.AppManifest = appManifestPathExpected
					logger.Debug(fmt.Sprintf("Found app manifest at: %s", ndkOptions.AppManifest))
				}
			}

			// Infer project root directory if not set
			if ndkOptions.ProjectRoot == "" {
				ndkOptions.ProjectRoot = filepath.Join(mergeNativeLibPath, "..", "..", "..", "..")
			}
		}

		// Build list of files to process based on whether path is a file or directory
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

	// Read AndroidManifest.xml if any key metadata is missing
	if ndkOptions.AppManifest != "" && (options.ApiKey == "" || ndkOptions.ApplicationId == "" || ndkOptions.VersionCode == "" || ndkOptions.VersionName == "") {
		logger.Debug("Reading data from AndroidManifest.xml")
		manifestData, err := android.ParseAndroidManifestXML(ndkOptions.AppManifest)
		if err != nil {
			return err
		}

		// Extract API key from manifest metadata if missing
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

		// Extract ApplicationId if missing
		if ndkOptions.ApplicationId == "" {
			ndkOptions.ApplicationId = manifestData.ApplicationId
			if ndkOptions.ApplicationId != "" {
				logger.Debug(fmt.Sprintf("Using %s as application ID from AndroidManifest.xml", ndkOptions.ApplicationId))
			}
		}

		// Extract VersionCode if missing
		if ndkOptions.VersionCode == "" {
			ndkOptions.VersionCode = manifestData.VersionCode
			if ndkOptions.VersionCode != "" {
				logger.Debug(fmt.Sprintf("Using %s as version code from AndroidManifest.xml", ndkOptions.VersionCode))
			}
		}

		// Extract VersionName if missing
		if ndkOptions.VersionName == "" {
			ndkOptions.VersionName = manifestData.VersionName
			if ndkOptions.VersionName != "" {
				logger.Debug(fmt.Sprintf("Using %s as version name from AndroidManifest.xml", ndkOptions.VersionName))
			}
		}
	}

	// Process native library files: convert .so to .so.sym via objcopy if needed
	for _, file := range fileList {
		if strings.HasSuffix(file, ".so.sym") {
			// Already a symbol file, add directly to upload list
			symbolFileList = append(symbolFileList, file)
		} else if soFileRegex.MatchString(file) {
			// Locate objcopy tool if not yet found
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

			// Create temp working directory for symbol extraction
			if workingDir == "" {
				workingDir, err = os.MkdirTemp("", "bugsnag-cli-ndk-*")
				if err != nil {
					return fmt.Errorf("error creating temporary working directory %s", err.Error())
				}
				defer os.RemoveAll(workingDir)
			}

			outputFile, err := android.Objcopy(objCopyPath, file, workingDir)
			if err != nil {
				return fmt.Errorf("failed to process file %s using objcopy: %s", file, err.Error())
			}

			symbolFileList = append(symbolFileList, outputFile)
		}
	}

	// Upload all extracted symbol files
	err = android.UploadAndroidNdk(
		symbolFileList,
		options.ApiKey,
		ndkOptions.ApplicationId,
		ndkOptions.VersionName,
		ndkOptions.VersionCode,
		ndkOptions.ProjectRoot,
		options,
		logger,
	)

	if err != nil {
		return err
	}

	return nil
}
