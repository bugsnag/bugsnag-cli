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

// DsymUploadInfo contains the relevant information for uploading dSYMs to bugsnag
type DsymUploadInfo struct {
	ProjectRoot string
	DsymPath    string
}

// XcodeBuildSettings contains the relevant build settings required for uploading to bugsnag
type XcodeBuildSettings struct {
	ConfigurationBuildDir string `mapstructure:"CONFIGURATION_BUILD_DIR"`
	InfoPlistPath         string `mapstructure:"INFOPLIST_PATH"`
	BuiltProductsDir      string `mapstructure:"BUILT_PRODUCTS_DIR"`
	DsymName              string `mapstructure:"DWARF_DSYM_FILE_NAME"`
}

// GetDefaultScheme checks if a scheme is in a given path or checks current directory if path is empty
func GetDefaultScheme(path, projectRoot string) (string, string, error) {
	schemes, derivedFrom := getXcodeSchemes(path, projectRoot)

	switch len(schemes) {
	case 0:
		return "", "", errors.Errorf("No schemes found in location '%s' please define which scheme to use with --scheme", path)
	case 1:
		return schemes[0], derivedFrom, nil
	default:
		return "", "", errors.Errorf("Multiple schemes found in location '%s', please define which scheme to use with --scheme", path)
	}
}

// IsSchemeInPath checks if a scheme is in a given path or checks current directory if path is empty
func IsSchemeInPath(path, schemeToFind, projectRoot string) (bool, string, error) {
	schemes, derivedFrom := getXcodeSchemes(path, projectRoot)
	for _, scheme := range schemes {
		if scheme == schemeToFind {
			return true, derivedFrom, nil
		}
	}

	return false, "", errors.Errorf("Unable to locate scheme '%s' in location: '%s'", schemeToFind, path)
}

// getXcodeSchemes parses the xcodebuild output for a given path to return a slice of schemes
func getXcodeSchemes(path, projectRoot string) ([]string, string) {
	var cmd *exec.Cmd

	// Used for showing which path value the scheme was found with i.e. either <path> or --project-root
	var derivedFrom string

	if isXcodebuildInstalled() {
		if strings.HasSuffix(path, ".xcworkspace") {
			cmd = exec.Command(utils.LocationOf(utils.XCODEBUILD), "-workspace", path, "-list")
			derivedFrom = path
		} else if strings.HasSuffix(path, ".xcodeproj") {
			cmd = exec.Command(utils.LocationOf(utils.XCODEBUILD), "-project", path, "-list")
			derivedFrom = path
		} else {

			cmd = exec.Command(utils.LocationOf(utils.XCODEBUILD), "-list")

			// Change the working directory of the command to path if projectRoot is not set, otherwise use projectRoot instead
			if projectRoot == "" {
				cmd.Dir = path
				derivedFrom = path
			} else {
				cmd.Dir = projectRoot
				derivedFrom = projectRoot
			}

		}
	} else {
		return []string{}, ""
	}

	output, err := cmd.Output()
	if err != nil {
		return []string{}, ""
	}

	schemes := strings.SplitAfterN(string(output), "Schemes:\n", 2)[1]

	// Remove excess whitespace and double newlines before splitting into a slice
	replacer := strings.NewReplacer(" ", "", "\n\n", "")
	sanitisedSchemes := replacer.Replace(schemes)

	schemesSlice := strings.Split(sanitisedSchemes, "\n")

	return schemesSlice, derivedFrom
}

// GetXcodeBuildSettings returns a struct of the relevant build settings for a given path and scheme
func GetXcodeBuildSettings(path, schemeName, projectRoot string) (*XcodeBuildSettings, error) {
	var buildSettings XcodeBuildSettings
	allBuildSettings, err := getXcodeBuildSettings(path, schemeName, projectRoot)
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
func getXcodeBuildSettings(path, schemeName, projectRoot string) (*map[string]*string, error) {
	var cmd *exec.Cmd

	if isXcodebuildInstalled() {
		if strings.HasSuffix(path, ".xcworkspace") {
			cmd = exec.Command(utils.LocationOf(utils.XCODEBUILD), "-workspace", path, "-scheme", schemeName, "-showBuildSettings")
		} else if strings.HasSuffix(path, ".xcodeproj") {
			cmd = exec.Command(utils.LocationOf(utils.XCODEBUILD), "-project", path, "-scheme", schemeName, "-showBuildSettings")
		} else {

			cmd = exec.Command(utils.LocationOf(utils.XCODEBUILD), "-scheme", schemeName, "-showBuildSettings")

			// Change the working directory of the command to path if projectRoot is not set, otherwise use projectRoot instead
			if projectRoot == "" {
				cmd.Dir = path
			} else {
				cmd.Dir = projectRoot
			}

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

// ProcessPathValue determines the projectRoot from a given path
func ProcessPathValue(path string, projectRoot string) (*DsymUploadInfo, error) {
	if path == "" && projectRoot == "" {
		currentDir, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		return &DsymUploadInfo{currentDir, ""}, nil
	}

	_, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	if utils.IsDir(path) {

		if projectRoot != "" {
			log.Info("--project-root flag set, it's value takes precedence and will be used for upload")
			return &DsymUploadInfo{projectRoot, ""}, nil
		}

		if strings.HasSuffix(path, ".xcodeproj") || strings.HasSuffix(path, ".xcworkspace") {
			// If path is pointing to a .xcodeproj or .xcworkspace directory, set projectRoot to one directory up
			return &DsymUploadInfo{filepath.Dir(path), ""}, nil
		} else {
			// If path is pointing to a directory, set projectRoot to the path
			return &DsymUploadInfo{path, ""}, nil
		}

	} else {
		// If path is pointing to a file, we will assume it's pointing to a dSYM and use as-is
		return &DsymUploadInfo{projectRoot, path}, nil
	}

}

// isXcodebuildInstalled checks if xcodebuild is installed by checking if there is a path returned for it
func isXcodebuildInstalled() bool {
	return utils.LocationOf(utils.XCODEBUILD) != ""
}
