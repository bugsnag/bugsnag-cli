package ios

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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
	DsymName              string `mapstructure:"DWARF_DSYM_FILE_NAME"`
	ProjectTempRoot       string `mapstructure:"PROJECT_TEMP_ROOT"`
}

// GetDefaultScheme checks if a scheme is in a given path or checks current directory if path is empty
func GetDefaultScheme(path string) (string, error) {
	schemes := getXcodeSchemes(path)

	switch len(schemes) {
	case 0:
		return "", errors.Errorf("No schemes found in location '%s' please define which scheme to use with --scheme", path)
	case 1:
		return schemes[0], nil
	default:
		return "", errors.Errorf("Multiple schemes found in location '%s', please define which scheme to use with --scheme", path)
	}
}

// IsSchemeInPath checks if a scheme is in a given path or checks current directory if path is empty
func IsSchemeInPath(path, schemeToFind string) (bool, error) {
	schemes := getXcodeSchemes(path)
	for _, scheme := range schemes {
		if scheme == schemeToFind {
			return true, nil
		}
	}

	return false, errors.Errorf("Unable to locate scheme '%s' in location: '%s'", schemeToFind, path)
}

// getXcodeSchemes parses the xcodebuild output for a given path to return a slice of schemes
func getXcodeSchemes(path string) []string {
	var cmd *exec.Cmd

	if isXcodebuildInstalled() {
		if strings.HasSuffix(path, ".xcworkspace") {
			cmd = exec.Command(utils.LocationOf(utils.XCODEBUILD), "-workspace", path, "-list")

		} else if strings.HasSuffix(path, ".xcodeproj") {
			cmd = exec.Command(utils.LocationOf(utils.XCODEBUILD), "-project", path, "-list")

		} else {

			cmd = exec.Command(utils.LocationOf(utils.XCODEBUILD), "-list")

			// Change the working directory of the command to path if it's a directory but not .xcodeproj or .xcworkspace
			cmd.Dir = path

		}
	} else {
		return []string{}
	}

	output, err := cmd.Output()
	if err != nil {
		return []string{}
	}

	schemes := strings.SplitAfterN(string(output), "Schemes:\n", 2)[1]

	// Remove excess whitespace and double newlines before splitting into a slice
	schemes = strings.ReplaceAll(schemes, "\n\n", "")
	schemesSlice := strings.Split(schemes, "\n")

	for i, scheme := range schemesSlice {
		schemesSlice[i] = strings.TrimSpace(scheme)
	}

	return schemesSlice
}

// GetXcodeBuildSettings returns a struct of the relevant build settings for a given path and scheme
func GetXcodeBuildSettings(path, schemeName string, configuration string) (*XcodeBuildSettings, error) {
	var buildSettings XcodeBuildSettings
	allBuildSettings, err := getXcodeBuildSettings(path, schemeName, configuration)
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
func getXcodeBuildSettings(path, schemeName string, configuration string) (*map[string]*string, error) {
	var cmd *exec.Cmd

	if isXcodebuildInstalled() {
		if !strings.HasSuffix(path, ".xcworkspace") && !strings.HasSuffix(path, ".xcodeproj") {
			path = FindXcodeProjOrWorkspace(path)
		}

		if strings.HasSuffix(path, ".xcworkspace") {
			cmdArgs := []string{"-workspace", path, "-scheme", schemeName, "-showBuildSettings"}
			if configuration != "" {
				cmdArgs = append(cmdArgs, "-configuration", configuration)
			}
			cmd = exec.Command(utils.LocationOf(utils.XCODEBUILD), cmdArgs...)
		} else if strings.HasSuffix(path, ".xcodeproj") {
			cmdArgs := []string{"-project", path, "-scheme", schemeName, "-showBuildSettings"}
			if configuration != "" {
				cmdArgs = append(cmdArgs, "-configuration", configuration)
			}
			cmd = exec.Command(utils.LocationOf(utils.XCODEBUILD), cmdArgs...)
		} else {
			return nil, fmt.Errorf("Unable to locate xcodeproj or xcworkspace in the given path")
		}
	} else {
		return nil, fmt.Errorf("Unable to locate xcodebuild on this system.")
	}

	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	buildSettings := strings.SplitAfterN(string(output), "Build settings for action build and target ", 2)[1]
	buildSettingsSlice := strings.Split(buildSettings, "\n")
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

func IsPathAnXcodeProjectOrWorkspace(path string) bool {
	if strings.HasSuffix(path, ".xcodeproj") || strings.HasSuffix(path, ".xcworkspace") {
		return true
	}

	var err error
	if isXcodebuildInstalled() {
		cmd := exec.Command(utils.LocationOf(utils.XCODEBUILD), "-list")
		cmd.Dir = path
		_, err = cmd.Output()
	} else {
		return false
	}

	return err == nil
}

// GetDefaultProjectRoot works out a value for using as project root if one isn't provided
func GetDefaultProjectRoot(path, projectRoot string) string {
	if projectRoot == "" {
		if path == "" {
			currentDir, _ := os.Getwd()
			return currentDir
		}

		if utils.IsDir(path) {

			// If path is pointing to a .xcodeproj or .xcworkspace directory, set the project root to one directory up
			if strings.HasSuffix(path, ".xcodeproj") || strings.HasSuffix(path, ".xcworkspace") {
				return filepath.Dir(path)

			}
		}

		// If path is pointing to a normal directory, set that as the project root
		return path

	} else {
		// If the project root is already set, use as-is
		return projectRoot
	}
}

// isXcodebuildInstalled checks if xcodebuild is installed by checking if there is a path returned for it
func isXcodebuildInstalled() bool {
	return utils.LocationOf(utils.XCODEBUILD) != ""
}

// FindXcodeProjOrWorkspace finds the .xcodeproj or .xcworkspace file in a given directory
// and returns the path to it
// If neither is found, an empty string is returned
// If both are found, the .xcworkspace file is returned
func FindXcodeProjOrWorkspace(path string) string {
	var xcodeProjPath string
	var xcodeWorkspacePath string

	files, err := os.ReadDir(path)
	if err != nil {
		return ""
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".xcodeproj") {
			xcodeProjPath = filepath.Join(path, file.Name())
		} else if strings.HasSuffix(file.Name(), ".xcworkspace") {
			xcodeWorkspacePath = filepath.Join(path, file.Name())
		}
	}

	if xcodeWorkspacePath != "" {
		return xcodeWorkspacePath
	} else if xcodeProjPath != "" {
		return xcodeProjPath
	}

	return ""
}
