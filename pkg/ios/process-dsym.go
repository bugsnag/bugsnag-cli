package ios

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/options"
	"github.com/bugsnag/bugsnag-cli/pkg/server"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
)

// ProcessDsymUpload uploads dSYM files to the specified endpoint.
// It retrieves the API key from the Info.plist if not provided, builds upload options,
// and uploads each dSYM file. If the initial upload fails with a 404 error, it retries at the base endpoint.
//
// Parameters:
// - plistPath: Path to the Info.plist file.
// - projectRoot: Root directory of the project.
// - options: CLI options containing configuration like API key.
// - dwarfInfo: List of dSYM information objects to be processed.
// - logger: Logger instance for debug and error messages.
//
// Returns:
// - An error if any part of the process fails; nil otherwise.
func ProcessDsymUpload(plistPath string, projectRoot string, options options.CLI, dwarfInfo []*DwarfInfo, logger log.Logger) error {
	var (
		plistData     *PlistData
		uploadOptions map[string]string
		err           error
	)

	// Retrieve API key from Info.plist if it exists and the API key is not already set.
	if utils.FileExists(plistPath) && options.ApiKey == "" {
		plistData, err = GetPlistData(plistPath)
		if err != nil {
			return fmt.Errorf("failed to read plist data: %w", err)
		}
		options.ApiKey = plistData.BugsnagProjectDetails.ApiKey
		if options.ApiKey != "" {
			logger.Debug(fmt.Sprintf("Using API key from Info.plist: %s", options.ApiKey))
		}
	}

	// Process and upload each dSYM file in the provided list.
	for _, dsym := range dwarfInfo {
		dsymInfo := fmt.Sprintf("(UUID: %s, Name: %s, Arch: %s)", dsym.UUID, dsym.Name, dsym.Arch)
		logger.Debug(fmt.Sprintf("Processing dSYM %s", dsymInfo))

		// Build upload options for the current dSYM file.
		uploadOptions, err = utils.BuildDsymUploadOptions(projectRoot)
		if err != nil {
			return fmt.Errorf("failed to build dSYM upload options: %w", err)
		}

		// Prepare the file data for uploading.
		fileFieldData := map[string]server.FileField{
			"dsym": server.LocalFile(filepath.Join(dsym.Location, dsym.Name)),
		}

		// Attempt to upload the dSYM file.
		err = server.ProcessFileRequest(options.ApiKey, "/dsym", uploadOptions, fileFieldData, dsym.UUID, options, logger)
		if err != nil {
			// Retry with the base endpoint if a 404 error occurs.
			if strings.Contains(err.Error(), "404 Not Found") {
				logger.Debug(fmt.Sprintf("Retrying upload for dSYM %s at base endpoint", dsymInfo))
				err = server.ProcessFileRequest(options.ApiKey, "", uploadOptions, fileFieldData, dsym.UUID, options, logger)
			}
			if err != nil {
				return fmt.Errorf("failed to upload dSYM %s: %w", dsymInfo, err)
			}
		}
	}

	return nil
}
