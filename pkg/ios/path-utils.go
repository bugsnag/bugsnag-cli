package ios

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"

	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
)

var DsymDirs []string
var TempDirs []string

// isPathADsymDirectory checks if the path is a directory containing dSYM(s) and makes a note of the dSYM locations
func isPathADsymDirectory(path string) bool {
	// If path is set and is a directory
	if path != "" && utils.IsDir(path) {
		// Check for dSYMs within it
		DsymDirs = append(DsymDirs, findDsyms(path)...)
		return len(DsymDirs) > 0

	} else {
		// If path is pointing to a .zip file, we will extract it and look for dSYMS within it to set DsymPaths
		if strings.HasSuffix(path, ".zip") {
			fileName := filepath.Base(path)
			log.Info("Attempting to unzip " + fileName + " before proceeding to upload")
			tempDir, err := utils.ExtractFile(path, "dsym")
			if err != nil {
				// TODO: This will be downgraded to a warning with --fail-on-upload in near future
				log.Error("Could not unzip "+fileName+" to a temporary directory, skipping", 1)
				// Silently remove the temp dir if one was created before continuing
				_ = os.RemoveAll(tempDir)
			} else {
				log.Info("Unzipped " + fileName + " to " + tempDir + " for uploading")
				TempDirs = append(TempDirs, tempDir)
				DsymDirs = append(DsymDirs, findDsyms(tempDir)...)
			}

			return len(TempDirs) > 0
		}

		// If path is pointing to a file with .dSYM suffix
		return strings.HasSuffix(path, ".dSYM")
	}

}

func ValidatePaths(path, dsymPath, projectRoot string) error {
	if isPathADsymDirectory(path) && projectRoot == "" || dsymPath != "" && projectRoot == "" {
		return errors.New("--project-root is mandatory when a path to dSYM(s) is supplied")
	}

	return nil
}

func SetWorkingDirectory(path string) string {
	if path == "" {
		currentDir, _ := os.Getwd()
		return currentDir
	}

	if utils.IsDir(path) {

		// If path is pointing to a .xcodeproj or .xcworkspace directory, set projectRoot to one directory up
		if strings.HasSuffix(path, ".xcodeproj") || strings.HasSuffix(path, ".xcworkspace") {
			// If path is pointing to a .xcodeproj or .xcworkspace directory, set projectRoot to one directory up
			return filepath.Dir(path)

		} else {

			if !isPathADsymDirectory(path) {
				return path
			} else {
				return ""
			}
		}

	} else {
		return ""
	}

}
