package ios

import (
	"os"

	"github.com/pkg/errors"
	"howett.net/plist"
)

// PlistData contains the relevant content of a plist file for uploading to Bugsnag.
// It extracts the app version, bundle version, and Bugsnag-specific project details.
type PlistData struct {
	BundleIdentifier      string                `plist:"CFBundleIdentifier"`
	VersionName           string                `plist:"CFBundleShortVersionString"`
	BundleVersion         string                `plist:"CFBundleVersion"`
	BugsnagProjectDetails bugsnagProjectDetails `plist:"bugsnag"`
}

type bugsnagProjectDetails struct {
	ApiKey string `plist:"apiKey"`
}

// GetPlistData parses the contents of an Info.plist file into a PlistData struct.
//
// This function reads the plist file directly and unmarshals its contents
// using the howett.net/plist library.
//
// Parameters:
// - plistFilePath: The path to the Info.plist file to be processed.
//
// Returns:
//   - A pointer to a PlistData struct containing parsed plist content.
//   - An error if the plist file cannot be processed or if required fields are missing.
func GetPlistData(plistFilePath string) (*PlistData, error) {
	if plistFilePath == "" {
		return nil, errors.New("plist file path is empty")
	}

	file, err := os.Open(plistFilePath)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open plist file: %s", plistFilePath)
	}
	defer file.Close()

	var plistData PlistData
	decoder := plist.NewDecoder(file)
	err = decoder.Decode(&plistData)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to parse plist data in file: %s", plistFilePath)
	}

	return &plistData, nil
}
