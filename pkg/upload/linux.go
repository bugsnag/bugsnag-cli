package upload

import (
	"fmt"
	"path/filepath"

	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/options"
	"github.com/bugsnag/bugsnag-cli/pkg/server"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
)

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
	if linuxOpts.VersionCode != "" {
		uploadOpts["versionCode"] = linuxOpts.VersionCode
	}
	if linuxOpts.ProjectRoot != "" {
		uploadOpts["projectRoot"] = linuxOpts.ProjectRoot
	}
	if base := filepath.Base(symbolFile); base != "" {
		uploadOpts["sharedObjectName"] = base
	}
	if linuxOpts.Overwrite {
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
//   - Resolves duplicates by ELF build ID.
//   - Reads metadata (appId, versionName) if provided.
//   - Uploads all recognized symbol files to the Bugsnag /linux endpoint.
//
// Returns:
//   - error: non-nil if scanning, build ID resolution, or upload fails.
func ProcessLinux(opts options.CLI, logger log.Logger) error {
	linuxOpts := opts.Upload.Linux

	var fileList []string
	var soFileList []string
	var err error

	for _, path := range linuxOpts.Path {
		logger.Info(fmt.Sprintf("Scanning path: %s", path))

		// Build a list of potential symbol files
		if utils.IsDir(path) {
			fileList, err = utils.BuildFileList([]string{path})
			if err != nil {
				return fmt.Errorf("building file list from %q: %w", path, err)
			}
			logger.Debug(fmt.Sprintf("Found %d files in directory %s", len(fileList), path))
		} else {
			fileList = append(fileList, path)
		}

		// Filter for valid ELF symbol files
		for _, file := range fileList {
			ok, err := utils.IsSymbolFile(file)
			if err != nil {
				logger.Debug(fmt.Sprintf("Skipping file %s", file))
				continue
			}
			if ok {
				soFileList = append(soFileList, file)
				logger.Debug(fmt.Sprintf("Found symbol file: %s", file))
			} else {
				logger.Debug(fmt.Sprintf("Skipping non-symbol file: %s", file))
			}
		}

		if linuxOpts.ProjectRoot != "" {
			logger.Debug(fmt.Sprintf("Using project root: %s", linuxOpts.ProjectRoot))
		}

		if len(soFileList) == 0 {
			logger.Warn("No symbol files found to upload")
			return nil
		}

		for _, file := range soFileList {
			if err := uploadSymbolFile(file, linuxOpts, opts, logger); err != nil {
				return err
			}
		}
	}
	return nil
}
