package ios

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"

	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
)

// XcodeBuildSettings contains the relevant build settings required for uploading to bugsnag
type XcodeBuildSettings struct {
	ConfigurationBuildDir string `mapstructure:"CONFIGURATION_BUILD_DIR"`
	InfoPlistPath         string `mapstructure:"INFOPLIST_PATH"`
	BuiltProductsDir      string `mapstructure:"BUILT_PRODUCTS_DIR"`
	DsymName              string `mapstructure:"DWARF_DSYM_FILE_NAME"`
}

// GetDefaultScheme checks if a scheme is in a given path or checks current directory if path is empty
func GetDefaultScheme(path string) (string, error) {
	schemes, err := getXcodeSchemes(path)
	if err != nil {
		return "", err
	}

	switch len(schemes) {
	case 0:
		return "", errors.Errorf("No schemes found in location '%s' please define which scheme to use with --scheme", path)
	case 1:
		return schemes[0], nil
	default:
		return "", errors.Errorf("No schemes found in location '%s', please define which scheme to use with --scheme", path)
	}
}

// IsSchemeInWorkspace checks if a scheme is in a given path or checks current directory if path is empty
func IsSchemeInWorkspace(path, schemeToFind string) (bool, error) {
	schemes, _ := getXcodeSchemes(path)
	for _, scheme := range schemes {
		if scheme == schemeToFind {
			return true, nil
		}
	}

	return false, errors.Errorf("Unable to locate scheme '%s' in location: '%s'", schemeToFind, path)
}

// getXcodeSchemes parses the xcodebuild output for a given path to return a slice of schemes
func getXcodeSchemes(path string) ([]string, error) {
	var cmd *exec.Cmd
	if isXcodebuildInstalled() {
		if strings.HasSuffix(path, ".xcworkspace") {
			cmd = exec.Command(utils.LocationOf(utils.XCODEBUILD), "-workspace", path, "-list")
		} else if strings.HasSuffix(path, ".xcodeproj") {
			cmd = exec.Command(utils.LocationOf(utils.XCODEBUILD), "-project", path, "-list")
		} else {
			cmd = exec.Command(utils.LocationOf(utils.XCODEBUILD), "-list")
		}
	} else {
		return nil, errors.New("Unable to locate xcodebuild on this system.")
	}

	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	schemes := strings.SplitAfterN(string(output), "Schemes:\n", 2)[1]

	// Remove excess whitespace and double newlines before splitting into a slice
	replacer := strings.NewReplacer(" ", "", "\n\n", "")
	sanitisedSchemes := replacer.Replace(schemes)

	schemesSlice := strings.Split(sanitisedSchemes, "\n")

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

// getXcodeBuildSettings parses the xcodebuild output for a given path and scheme to return a map of all build settings
func getXcodeBuildSettings(path, schemeName string) (*map[string]*string, error) {
	var cmd *exec.Cmd
	if isXcodebuildInstalled() {
		if strings.HasSuffix(path, ".xcworkspace") {
			cmd = exec.Command(utils.LocationOf(utils.XCODEBUILD), "-workspace", path, "-scheme", schemeName, "-showBuildSettings")
		} else if strings.HasSuffix(path, ".xcodeproj") {
			cmd = exec.Command(utils.LocationOf(utils.XCODEBUILD), "-project", path, "-scheme", schemeName, "-showBuildSettings")
		} else {
			cmd = exec.Command(utils.LocationOf(utils.XCODEBUILD), "-scheme", schemeName, "-showBuildSettings")
		}
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

// GetProjectRoot determines the projectRoot from a given path
func GetProjectRoot(path string, projRootSet bool) (string, error) {
	var projectRoot string

	if projRootSet {
		log.Info("--project-root flag set, it's value takes precedence and will be used for upload")
		return path, nil
	}

	_, err := os.Stat(path)
	if err != nil {
		return "", err
	}

	if utils.IsDir(path) {

		if strings.HasSuffix(path, ".xcworkspace") {
			// If path is pointing to a .xcworkspace directory, set projectRoot to two directory up
			projectRoot = filepath.Dir(filepath.Dir(path))

		} else if strings.HasSuffix(path, ".xcodeproj") {
			// If path is pointing to a .xcworkspace directory, set projectRoot to one directory up
			projectRoot = filepath.Dir(path)
		}

	} else {
		log.Error("string argument passed to GetProjectRoot is not a directory", 1)
	}

	return "", err
}

// isXcodebuildInstalled checks if xcodebuild is installed by checking if there is a path returned for it
func isXcodebuildInstalled() bool {
	return utils.LocationOf(utils.XCODEBUILD) != ""
}
