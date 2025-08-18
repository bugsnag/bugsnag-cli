package upload

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/options"
	"github.com/bugsnag/bugsnag-cli/pkg/server"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
)

// buildVegaVariantFolders scans a directory and returns a list of folders
// matching the pattern private/vega/{arch}/{variant}/lib/{arch}/.
func buildVegaVariantFolders(baseDir, variant string) ([]string, error) {
	var matches []string

	vegaPath := filepath.Join(baseDir, "private", "vega")

	entries, err := os.ReadDir(vegaPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read vega directory: %w", err)
	}

	// First level: arch
	for _, archEntry := range entries {
		if !archEntry.IsDir() {
			continue
		}
		arch := archEntry.Name()
		variantPath := filepath.Join(vegaPath, arch, variant, "lib", arch)

		if utils.IsDir(variantPath) {
			matches = append(matches, variantPath)
		}
	}

	return matches, nil
}

// buildArchVariantFolders scans a directory and returns a list of subfolders
// that match the format "arch-variant" (e.g. "arm64-debug", "x86-release").
func buildArchVariantFolders(baseDir string, variant string) ([]string, error) {
	var matches []string

	// Regex: arch-variant (letters, numbers, underscores, dashes allowed), with specific variant
	re := regexp.MustCompile(`^[a-zA-Z0-9_]+-` + regexp.QuoteMeta(variant) + `$`)

	entries, err := os.ReadDir(baseDir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() && re.MatchString(entry.Name()) {
			matches = append(matches, filepath.Join(baseDir, entry.Name(), "debug"))
		}
	}

	return matches, nil
}

func ProcessLinux(opts options.CLI, logger log.Logger) error {
	var (
		err              error
		symbolFiles      []string
		symbolFolder     []string
		vegaSymbolFolder []string
		vegaSymbolFiles  []string
	)

	linuxOptions := opts.Upload.Linux

	for _, path := range linuxOptions.Path {
		// Locate the build folder
		if linuxOptions.BuildFolder == "" {
			possibleBuildFolder := filepath.Join(path, "build")
			if utils.IsDir(possibleBuildFolder) {
				linuxOptions.BuildFolder = utils.Path(possibleBuildFolder)
			} else {
				return fmt.Errorf("unable to find the build directory in %s", path)
			}
		}

		logger.Info(fmt.Sprintf("Searching for symbol files in %s for variant %s", linuxOptions.BuildFolder, linuxOptions.Variant))

		//	Locate symbol files in build folder
		//	Paths to check
		//	- build/{arch}-{variant}/debug/*.so.debug
		//	- build/private/vega/{arch}/{variant}/lib/{arch}/*.so

		//	Build a list of folders for symbolFiles
		symbolFolder, err = buildArchVariantFolders(string(linuxOptions.BuildFolder), linuxOptions.Variant)

		if err != nil {
			return err
		}

		symbolFiles, err = utils.BuildFileList(symbolFolder)

		if err != nil {
			return err
		}

		logger.Info(fmt.Sprintf("Found %d symbol files", len(symbolFiles)))

		logger.Info(fmt.Sprintf("Searching for vega symbol files in %s for variant %s", linuxOptions.BuildFolder, utils.Capitalize(linuxOptions.Variant)))

		for _, file := range symbolFiles {
			logger.Info(fmt.Sprintf("Uploading %s", file))
		}

		vegaSymbolFolder, err = buildVegaVariantFolders(string(linuxOptions.BuildFolder), utils.Capitalize(linuxOptions.Variant))

		if err != nil {
			return err
		}

		vegaSymbolFiles, err = utils.BuildFileList(vegaSymbolFolder)

		logger.Info(fmt.Sprintf("Found %d vega symbol files", len(vegaSymbolFiles)))

		for _, file := range vegaSymbolFiles {
			logger.Info(fmt.Sprintf("Uploading %s", file))
		}

		logger.Info("Building Upload Options")

		for _, symbolFile := range symbolFiles {
			uploadOpts := map[string]string{}

			// Can be gotten from the manifest.toml file
			if linuxOptions.ApplicationId != "" {
				uploadOpts["appId"] = linuxOptions.ApplicationId
			}
			// Can be gotten from the manifest.toml file
			if linuxOptions.VersionCode != "" {
				uploadOpts["versionCode"] = linuxOptions.VersionCode
			}
			// Can be gotten from the manifest.toml file
			if linuxOptions.VersionName != "" {
				uploadOpts["versionName"] = linuxOptions.VersionName
			}
			if linuxOptions.ProjectRoot != "" {
				uploadOpts["projectRoot"] = linuxOptions.ProjectRoot
			}
			if base := filepath.Base(symbolFile); base != "" {
				uploadOpts["sharedObjectName"] = base
			}
			if opts.Upload.Overwrite {
				uploadOpts["overwrite"] = "true"
			}

			fileField := map[string]server.FileField{
				"soFile": server.LocalFile(symbolFile),
			}

			err := server.ProcessFileRequest(
				opts.ApiKey,
				"/linux",
				uploadOpts,
				fileField,
				filepath.Base(symbolFile),
				opts,
				logger,
			)
			if err != nil {
				return fmt.Errorf("failed to upload NDK symbol for %s: %w", symbolFile, err)
			}

		}

	}

	return nil
}
