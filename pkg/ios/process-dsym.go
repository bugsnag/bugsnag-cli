package ios

import (
	"fmt"
	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/options"
	"github.com/bugsnag/bugsnag-cli/pkg/server"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
	"path/filepath"
	"strings"
)

func ProcessDsymUpload(plistPath, endpoint, projectRoot string, options options.CLI, logger log.Logger) error {
	var plistData *PlistData
	var dwarfInfo []*DwarfInfo
	var uploadOptions map[string]string

	var err error

	if utils.FileExists(plistPath) && options.ApiKey == "" {
		plistData, err = GetPlistData(plistPath)
		if err != nil {
			return err
		}
		options.ApiKey = plistData.BugsnagProjectDetails.ApiKey
		if options.ApiKey != "" {
			logger.Debug(fmt.Sprintf("Using API key from Info.plist: %s", options.ApiKey))
		}
	}

	for _, dsym := range dwarfInfo {
		dsymInfo := fmt.Sprintf("(UUID: %s, Name: %s, Arch: %s)", dsym.UUID, dsym.Name, dsym.Arch)
		logger.Debug(fmt.Sprintf("Processing dSYM %s", dsymInfo))

		// Build upload options for each dSYM file
		uploadOptions, err = utils.BuildDsymUploadOptions(options.ApiKey, projectRoot)
		if err != nil {
			return err
		}

		// Prepare the file data for uploading
		fileFieldData := map[string]server.FileField{
			"dsym": server.LocalFile(filepath.Join(dsym.Location, dsym.Name)),
		}

		// Attempt to upload the dSYM file to the endpoint
		err = server.ProcessFileRequest(endpoint+"/dsym", uploadOptions, fileFieldData, dsym.UUID, options, logger)
		if err != nil && strings.Contains(err.Error(), "404 Not Found") {
			// If the first upload fails due to 404, retry uploading to the base endpoint
			err = server.ProcessFileRequest(endpoint, uploadOptions, fileFieldData, dsym.UUID, options, logger)
		}
		if err != nil {
			return err
		}
	}
	return nil
}
