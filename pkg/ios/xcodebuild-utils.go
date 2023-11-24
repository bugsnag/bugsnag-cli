package ios

import (
	"os/exec"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"

	"github.com/bugsnag/bugsnag-cli/pkg/utils"
)

// XcodeBuildSettings contains the relevant build settings required for uploading to bugsnag
type XcodeBuildSettings struct {
	ConfigurationBuildDir string `mapstructure:"CONFIGURATION_BUILD_DIR"`
	InfoPlistPath         string `mapstructure:"INFOPLIST_PATH"`
	BuiltProductsDir      string `mapstructure:"BUILT_PRODUCTS_DIR"`
}

// IsSchemeInWorkspace checks if a scheme is in a given workspace
func IsSchemeInWorkspace(workspacePath, schemeToFind string) (bool, error) {
	schemes, err := getXcodeSchemes(workspacePath)
	if err != nil {
		return false, err
	}

	for _, scheme := range schemes {
		if scheme == schemeToFind {
			return true, nil
		}
	}

	return false, errors.New("Unable to locate a suitable .xcworkspace file")
}

// getXcodeSchemes parses the xcodebuild output for a given workspace path to return a slice of schemes
func getXcodeSchemes(workspacePath string) ([]string, error) {
	var cmd *exec.Cmd
	if isXcodebuildInstalled() {
		cmd = exec.Command("xcodebuild", "-workspace", workspacePath, "-list")
	} else {
		return nil, errors.New("Unable to locate xcodebuild on this system.")
	}

	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	schemes := strings.SplitAfterN(string(output), "Schemes:", 2)[1]
	schemesSlice := strings.Split(strings.ReplaceAll(schemes, " ", ""), "\n")

	return schemesSlice, nil
}

// GetXcodeBuildSettings returns a struct of the relevant build settings for a given workspace and scheme
func GetXcodeBuildSettings(workspacePath, schemeName string) (*XcodeBuildSettings, error) {
	var buildSettings XcodeBuildSettings
	allBuildSettings, err := getXcodeBuildSettings(workspacePath, schemeName)
	if err != nil {
		return nil, err
	}
	err = mapstructure.Decode(allBuildSettings, &buildSettings)
	if err != nil {
		return nil, err
	}

	return &buildSettings, nil
}

// getXcodeBuildSettings parses the xcodebuild output for a given workspace and scheme to return a map of all build settings
func getXcodeBuildSettings(workspacePath, schemeName string) (*map[string]*string, error) {
	var cmd *exec.Cmd
	if isXcodebuildInstalled() {
		cmd = exec.Command("xcodebuild", "-workspace", workspacePath, "-scheme", schemeName, "-showBuildSettings")
	} else {
		return nil, errors.New("Unable to locate xcodebuild on this system.")
	}

	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	buildSettings := strings.SplitAfterN(string(output), "Build settings for action build and target ", 2)[1]
	buildSettingsSlice := strings.Split(strings.ReplaceAll(buildSettings, " ", ""), "\n")

	buildSettingsMap := make(map[string]*string)
	for _, line := range buildSettingsSlice {
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			buildSettingsMap[key] = &value
		}
	}

	return &buildSettingsMap, nil
}

// isXcodebuildInstalled checks if xcodebuild is installed by checking if there is a path returned for it
func isXcodebuildInstalled() bool {
	return utils.FindLocationOf(utils.XCODEBUILD) != ""
}
