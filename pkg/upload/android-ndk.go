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

// resolveMergedLibPath resolves the absolute path to the merged_native_libs directory
// for the given input file or directory.
//
// Parameters:
//   - input: path to a .so file or a subdirectory of an Android project.
//
// Returns:
//   - string: the resolved merged_native_libs path.
//   - error: non-nil if resolution fails.
func resolveMergedLibPath(input string) (string, error) {
	return android.FindNativeLibPath([]string{"android", "app", "build", "intermediates", "merged_native_libs"}, input)
}

// resolveAppManifestIfNeeded sets the AndroidManifest.xml path in ndkOpts if it hasn't already been set.
//
// It infers the manifest location based on the provided native lib path and variant.
//
// Parameters:
//   - ndkOpts: AndroidNdkMapping options struct (will be mutated).
//   - libPath: resolved path to merged_native_libs.
//   - logger: logger used to emit debug output.
func resolveAppManifestIfNeeded(ndkOpts *options.AndroidNdkMapping, libPath string, logger log.Logger) {
	if ndkOpts.AppManifest != "" {
		return
	}
	manifestPath := filepath.Join(libPath, "..", "merged_manifests", ndkOpts.Variant, "AndroidManifest.xml")
	if utils.FileExists(manifestPath) {
		ndkOpts.AppManifest = manifestPath
		logger.Debug(fmt.Sprintf("Found AndroidManifest.xml at %s", ndkOpts.AppManifest))
	}
}

// resolveProjectRootIfNeeded sets the project root directory in ndkOpts if it hasn't already been set.
//
// It infers the root by navigating upward from the provided native lib path.
//
// Parameters:
//   - ndkOpts: AndroidNdkMapping options struct (will be mutated).
//   - libPath: resolved path to merged_native_libs.
func resolveProjectRootIfNeeded(ndkOpts *options.AndroidNdkMapping, libPath string) {
	if ndkOpts.ProjectRoot == "" {
		ndkOpts.ProjectRoot = filepath.Join(libPath, "..", "..", "..", "..")
	}
}

// resolveFileList generates a list of .so or .so.sym files to be processed from the given input path.
//
// If the input path is a file, it returns a single-element list. If it's a directory,
// it resolves files from either the path directly or from a subdirectory matching the variant.
//
// Parameters:
//   - inputPath: user-supplied path to a .so file or directory.
//   - mergedLibPath: base path to merged_native_libs.
//   - variant: build variant (e.g. "release", "debug").
//
// Returns:
//   - []string: list of resolved file paths.
//   - error: non-nil if resolution fails.
func resolveFileList(inputPath, mergedLibPath, variant string) ([]string, error) {
	if !utils.IsDir(inputPath) {
		return []string{inputPath}, nil
	}
	if strings.Contains(inputPath, filepath.Join("merged_native_libs", variant)) {
		return utils.BuildFileList([]string{inputPath})
	}
	return utils.BuildFileList([]string{filepath.Join(mergedLibPath, variant)})
}

// populateMetadataFromManifest extracts Bugsnag metadata from AndroidManifest.xml and populates CLI options.
//
// Metadata includes the API key, application ID, version name, and version code if not already set.
//
// Parameters:
//   - opts: pointer to CLI options struct (may be mutated).
//   - ndkOpts: pointer to AndroidNdkMapping struct (may be mutated).
//   - logger: logger used to emit debug output.
//
// Returns:
//   - error: non-nil if manifest parsing fails.
func populateMetadataFromManifest(opts *options.CLI, ndkOpts *options.AndroidNdkMapping, logger log.Logger) error {
	logger.Debug("Parsing metadata from AndroidManifest.xml")
	manifestData, err := android.ParseAndroidManifestXML(ndkOpts.AppManifest)
	if err != nil {
		return err
	}

	if opts.ApiKey == "" {
		for key, name := range manifestData.Application.MetaData.Name {
			if name == "com.bugsnag.android.API_KEY" {
				opts.ApiKey = manifestData.Application.MetaData.Value[key]
				logger.Debug(fmt.Sprintf("API key found: %s", opts.ApiKey))
				break
			}
		}
	}

	if ndkOpts.ApplicationId == "" {
		ndkOpts.ApplicationId = manifestData.ApplicationId
		logger.Debug(fmt.Sprintf("ApplicationId: %s", ndkOpts.ApplicationId))
	}

	if ndkOpts.VersionCode == "" {
		ndkOpts.VersionCode = manifestData.VersionCode
		logger.Debug(fmt.Sprintf("VersionCode: %s", ndkOpts.VersionCode))
	}

	if ndkOpts.VersionName == "" {
		ndkOpts.VersionName = manifestData.VersionName
		logger.Debug(fmt.Sprintf("VersionName: %s", ndkOpts.VersionName))
	}

	return nil
}

// ProcessAndroidNDK processes Android NDK symbol files for uploading to Bugsnag.
//
// It performs the following steps:
//   - Resolves native libraries and variant information from project paths.
//   - Parses metadata from AndroidManifest.xml if needed.
//   - Extracts debug symbols from .so files using objcopy.
//   - Uploads the resulting .so.sym files and metadata to Bugsnag.
//
// Parameters:
//   - opts: CLI options including upload config and metadata.
//   - logger: logger used to emit debug output.
//
// Returns:
//   - error: non-nil if any processing or upload step fails.
func ProcessAndroidNDK(opts options.CLI, logger log.Logger) error {
	ndkOpts := opts.Upload.AndroidNdk
	soRegex := regexp.MustCompile(`\.so.*$`)

	var (
		fileList    []string
		objCopyPath string
		workingDir  string
		err         error
	)

	for _, inputPath := range ndkOpts.Path {
		libPath, err := resolveMergedLibPath(inputPath)
		if err != nil {
			return err
		}

		if filepath.Base(libPath) == "merged_native_libs" {
			if ndkOpts.Variant == "" {
				ndkOpts.Variant, err = android.GetVariantDirectory(libPath)
				if err != nil {
					return err
				}
			}
			resolveAppManifestIfNeeded(&ndkOpts, libPath, logger)
			resolveProjectRootIfNeeded(&ndkOpts, libPath)
		}

		files, err := resolveFileList(inputPath, libPath, ndkOpts.Variant)
		if err != nil {
			return fmt.Errorf("building file list for variant %q: %w", ndkOpts.Variant, err)
		}
		fileList = append(fileList, files...)
	}

	if ndkOpts.AppManifest != "" && (opts.ApiKey == "" || ndkOpts.ApplicationId == "" || ndkOpts.VersionCode == "" || ndkOpts.VersionName == "") {
		if err := populateMetadataFromManifest(&opts, &ndkOpts, logger); err != nil {
			return err
		}
	}

	if ndkOpts.ProjectRoot != "" {
		logger.Debug(fmt.Sprintf("Using %s as the project root", ndkOpts.ProjectRoot))
	}

	for _, file := range fileList {
		symbols := make(map[string]string)

		if strings.HasSuffix(file, ".so.sym") {
			symbols[file] = file
		} else if soRegex.MatchString(file) {
			if objCopyPath == "" {
				ndkOpts.AndroidNdkRoot, err = android.GetAndroidNDKRoot(ndkOpts.AndroidNdkRoot)
				if err != nil {
					return err
				}
				objCopyPath, err = android.BuildObjcopyPath(ndkOpts.AndroidNdkRoot)
				if err != nil {
					return err
				}
				logger.Debug(fmt.Sprintf("Using objcopy from NDK: %s", objCopyPath))
			}

			if workingDir == "" {
				workingDir, err = os.MkdirTemp("", "bugsnag-cli-ndk-*")
				if err != nil {
					return fmt.Errorf("creating temp directory: %w", err)
				}
				defer os.RemoveAll(workingDir)
			}

			logger.Debug(fmt.Sprintf("Extracting symbols from %s", file))
			outputFile, err := android.Objcopy(objCopyPath, file, workingDir)
			if err != nil {
				return fmt.Errorf("objcopy failed for %s: %w", file, err)
			}
			logger.Debug(fmt.Sprintf("Extracted symbol files to %s", outputFile))
			symbols[file] = outputFile
		}

		if err := android.UploadAndroidNdk(
			symbols,
			opts.ApiKey,
			ndkOpts.ApplicationId,
			ndkOpts.VersionName,
			ndkOpts.VersionCode,
			ndkOpts.ProjectRoot,
			opts,
			logger,
		); err != nil {
			return err
		}
	}

	return nil
}
