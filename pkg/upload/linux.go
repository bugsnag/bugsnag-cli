package upload

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/bugsnag/bugsnag-cli/pkg/linux"
	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/options"
	"github.com/bugsnag/bugsnag-cli/pkg/server"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
)

// buildArchVariantFolders scans a build directory and returns subfolders that
// match the naming convention "arch-variant" (e.g. "arm64-debug", "x86-release").
//
// Parameters:
//   - baseDir: The root directory containing arch-variant subfolders.
//   - variant: The variant string to match (e.g. "debug" or "release").
//
// Returns:
//   - []string: A list of matching directories pointing to their "debug" subfolder.
//   - error: non-nil if the baseDir cannot be read.
func buildArchVariantFolders(baseDir, variant string) ([]string, error) {
	var matches []string

	re := regexp.MustCompile(`^[a-zA-Z0-9_]+-` + regexp.QuoteMeta(variant) + `$`)

	entries, err := os.ReadDir(baseDir)
	if err != nil {
		return nil, fmt.Errorf("reading directory %q: %w", baseDir, err)
	}

	for _, entry := range entries {
		if entry.IsDir() && re.MatchString(entry.Name()) {
			matches = append(matches, filepath.Join(baseDir, entry.Name(), "debug"))
		}
	}

	return matches, nil
}

// uploadSymbolFile uploads a single Linux symbol file to the Bugsnag symbol server.
//
// Parameters:
//   - symbolFile: The path to the symbol file to upload.
//   - linuxOpts: Linux-specific upload options including appId, versionName, etc.
//   - opts: Global CLI options including an API key and overwrite behavior.
//   - logger: Logger for structured output.
//
// Returns:
//   - error: non-nil if the upload fails due to request or file issues.
func uploadSymbolFile(symbolFile string, linuxOpts options.LinuxOptions, opts options.CLI, logger log.Logger) error {
	uploadOpts := map[string]string{}

	if linuxOpts.ApplicationId != "" {
		uploadOpts["appId"] = linuxOpts.ApplicationId
	}
	if linuxOpts.VersionName != "" {
		uploadOpts["versionName"] = linuxOpts.VersionName
	}
	if linuxOpts.ProjectRoot != "" {
		uploadOpts["projectRoot"] = linuxOpts.ProjectRoot
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

	if err := server.ProcessFileRequest(
		opts.ApiKey,
		"/linux",
		uploadOpts,
		fileField,
		filepath.Base(symbolFile),
		opts,
		logger,
	); err != nil {
		return fmt.Errorf("uploading Linux symbol file %q: %w", symbolFile, err)
	}
	return nil
}

// ProcessLinux locates, validates, and uploads Linux symbol files.
//
// Parameters:
//   - opts: Global CLI options including upload configuration and API key.
//   - logger: Logger for structured logging and debug output.
//
// Behavior:
//   - Scans provided paths for build folders or symbol files.
//   - Resolves arch-variant directories when a build folder is found.
//   - Reads metadata (appId, versionName) from manifest.toml if missing.
//   - Uploads all recognized symbol files to the Bugsnag /linux endpoint.
//
// Returns:
//   - error: non-nil if scanning, manifest reading, or upload fails.
func ProcessLinux(opts options.CLI, logger log.Logger) error {
	linuxOpts := opts.Upload.Linux

	for _, path := range linuxOpts.Path {
		var (
			searchDir   string
			symbolDirs  []string
			symbolFiles []string
		)

		if utils.IsDir(path) {
			// Determine the correct folder to search
			switch {
			case linuxOpts.BuildFolder != "":
				searchDir = string(linuxOpts.BuildFolder)
			case utils.IsDir(filepath.Join(path, "build")):
				searchDir = filepath.Join(path, "build")
			case filepath.Base(path) == "build":
				searchDir = path
			default:
				// Directly treat as a folder containing symbol files
				symbolDirs = []string{path}
			}

			// If a build directory is set, scan arch-variant folders
			if searchDir != "" {
				logger.Info(fmt.Sprintf("Scanning %s for symbol files (variant: %s)", searchDir, linuxOpts.Variant))
				var err error
				symbolDirs, err = buildArchVariantFolders(searchDir, linuxOpts.Variant)
				if err != nil {
					return fmt.Errorf("scanning arch-variant folders in %q: %w", searchDir, err)
				}
			}

			var err error
			symbolFiles, err = utils.BuildFileList(symbolDirs)
			if err != nil {
				return fmt.Errorf("building symbol file list from %v: %w", symbolDirs, err)
			}
		} else {
			// Single file provided
			ok, err := utils.IsSymbolFile(path)
			if err != nil {
				return fmt.Errorf("checking if %q is a symbol file: %w", path, err)
			}
			if ok {
				symbolFiles, err = utils.BuildFileList([]string{path})
				if err != nil {
					return fmt.Errorf("building symbol file list from single file %q: %w", path, err)
				}
			}
		}

		logger.Info(fmt.Sprintf("Found %d symbol files to upload", len(symbolFiles)))

		// If metadata missing, attempt to read from manifest.toml
		if linuxOpts.ApplicationId == "" || linuxOpts.VersionName == "" {
			manifestPath := filepath.Join(path, "manifest.toml")
			logger.Debug(fmt.Sprintf("Attempting to read metadata from %s", manifestPath))
			version, appID, err := linux.ReadManifest(manifestPath)
			if err != nil {
				logger.Warn(fmt.Sprintf("Unable to read manifest.toml at %s: %s", manifestPath, err.Error()))
			} else {
				linuxOpts.VersionName = version
				linuxOpts.ApplicationId = appID
			}
		}

		if linuxOpts.ProjectRoot != "" {
			logger.Debug(fmt.Sprintf("Using %s as project root", linuxOpts.ProjectRoot))
		}

		for _, file := range symbolFiles {
			if err := uploadSymbolFile(file, linuxOpts, opts, logger); err != nil {
				return err
			}
		}
	}

	return nil
}
