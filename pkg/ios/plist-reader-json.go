package ios

import (
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/pkg/errors"

	"github.com/bugsnag/bugsnag-cli/pkg/utils"
)

// PlistData contains the relevant content of a plist file for uploading to Bugsnag.
// It extracts the app version, bundle version, and Bugsnag-specific project details.
type PlistData struct {
	VersionName           string                `json:"CFBundleShortVersionString"`
	BundleVersion         string                `json:"CFBundleVersion"`
	BugsnagProjectDetails bugsnagProjectDetails `json:"bugsnag"`
}

type bugsnagProjectDetails struct {
	ApiKey string `json:"apiKey"`
}

// GetPlistData parses the contents of an Info.plist file into a PlistData struct.
//
// This function uses the `plutil` command-line utility to convert the plist file
// into a JSON representation, which is then unmarshaled into a Go struct.
//
// Parameters:
// - plistFilePath: The path to the Info.plist file to be processed.
//
// Returns:
//   - A pointer to a PlistData struct containing parsed plist content.
//   - An error if the plist file cannot be processed, if `plutil` is unavailable,
//     or if required fields are missing in the plist data.
func GetPlistData(plistFilePath string) (*PlistData, error) {
	if plistFilePath == "" {
		return nil, errors.New("plist file path is empty")
	}

	if !isPlutilInstalled() {
		return nil, errors.New("plutil is not installed or could not be located")
	}

	cmd := exec.Command(utils.LocationOf(utils.PLUTIL), "-convert", "json", "-o", "-", plistFilePath)

	output, err := cmd.Output()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to execute plutil on file: %s", plistFilePath)
	}

	if len(output) == 0 {
		return nil, errors.New(fmt.Sprintf("plutil returned empty output reading file: %s", plistFilePath))
	}

	var plistData PlistData
	err = json.Unmarshal(output, &plistData)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("unable to parse plist data in file: %s", plistFilePath))
	}

	// Validate required fields in the plist data
	if plistData.VersionName == "" || plistData.BundleVersion == "" {
		return nil, fmt.Errorf("invalid plist data: missing required fields (VersionName or BundleVersion)")
	}

	return &plistData, nil
}

// isPlutilInstalled checks if the `plutil` utility is installed on the system.
//
// This function verifies the availability of `plutil` by checking its system path
// using a utility function.
//
// Returns:
// - `true` if `plutil` is found in the system's executable path.
// - `false` otherwise.
func isPlutilInstalled() bool {
	return utils.LocationOf(utils.PLUTIL) != ""
}
