package ios

import (
	"encoding/json"
	"os/exec"

	"github.com/pkg/errors"

	"github.com/bugsnag/bugsnag-cli/pkg/utils"
)

// PlistData contains the relevant content of a plist file for uploading to bugsnag
type PlistData struct {
	VersionName           string                `json:"CFBundleShortVersionString"`
	BundleVersion         string                `json:"CFBundleVersion"`
	BugsnagProjectDetails bugsnagProjectDetails `json:"bugsnag"`
}

type bugsnagProjectDetails struct {
	ApiKey string `json:"apiKey"`
}

// GetPlistData returns the relevant content of a plist file as a PlistData struct
func GetPlistData(plistFilePath string) (*PlistData, error) {
	var plistData *PlistData
	var cmd *exec.Cmd

	if isPlutilInstalled() {
		cmd = exec.Command(utils.FindLocationOf(utils.PLUTIL), "-convert", "json", "-o", "-", plistFilePath)

		output, err := cmd.Output()
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(output, &plistData)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("Unable to locate plutil in it's default location `/usr/bin/plutil` on this system.")
	}

	return plistData, nil
}

// isPlutilInstalled checks if plutil is installed by checking if there is a path returned for it
func isPlutilInstalled() bool {
	return utils.FindLocationOf(utils.PLUTIL) != ""
}
