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

// XcodeBuildSettings contains the relevant build settings required for uploading to Bugsnag.
type XcodeBuildSettings struct {
	ConfigurationBuildDir string `mapstructure:"CONFIGURATION_BUILD_DIR"`
	InfoPlistPath         string `mapstructure:"INFOPLIST_PATH"`
	BuiltProductsDir      string `mapstructure:"BUILT_PRODUCTS_DIR"`
	DsymName              string `mapstructure:"DWARF_DSYM_FILE_NAME"`
	ProjectTempRoot       string `mapstructure:"PROJECT_TEMP_ROOT"`
}

// GetDefaultScheme determines the default Xcode scheme in a given path or the current directory if no path is provided.
//
// Parameters:
// - path (string): Path to search for schemes.
//
// Returns:
// - string: The name of the default scheme.
// - error: If no schemes or multiple schemes are found.
func GetDefaultScheme(path string) (string, error) {
	schemes := getXcodeSchemes(path)

	switch len(schemes) {
	case 0:
		return "", errors.Errorf("no schemes found in location '%s'. Please specify a scheme with --scheme", path)
	case 1:
		return schemes[0], nil
	default:
		return "", errors.Errorf("multiple schemes found in location '%s'. Please specify a scheme with --scheme", path)
	}
}

// IsSchemeInPath verifies whether a given scheme exists in a specified path or current directory.
//
// Parameters:
// - path (string): Path to search for the scheme.
// - schemeToFind (string): Scheme name to look for.
//
// Returns:
// - bool: True if the scheme exists; false otherwise.
// - error: If the scheme cannot be located.
func IsSchemeInPath(path, schemeToFind string) (bool, error) {
	schemes := getXcodeSchemes(path)
	for _, scheme := range schemes {
		if scheme == schemeToFind {
			return true, nil
		}
	}
	return false, errors.Errorf("unable to locate scheme '%s' in location '%s'", schemeToFind, path)
}

// getXcodeSchemes retrieves a list of Xcode schemes by parsing the `xcodebuild` output.
//
// Parameters:
// - path (string): Path to search for schemes.
//
// Returns:
// - []string: A slice of scheme names.
func getXcodeSchemes(path string) []string {
	var cmd *exec.Cmd

	if isXcodebuildInstalled() {
		if strings.HasSuffix(path, ".xcworkspace") {
			cmd = exec.Command(utils.LocationOf(utils.XCODEBUILD), "-workspace", path, "-list")
		} else if strings.HasSuffix(path, ".xcodeproj") {
			cmd = exec.Command(utils.LocationOf(utils.XCODEBUILD), "-project", path, "-list")
		} else {
			cmd = exec.Command(utils.LocationOf(utils.XCODEBUILD), "-list")
			cmd.Dir = path // Set working directory if path is a directory
		}
	} else {
		return []string{}
	}

	output, err := cmd.Output()
	if err != nil {
		return []string{}
	}

	schemes := strings.SplitAfterN(string(output), "Schemes:\n", 2)[1]
	schemes = strings.ReplaceAll(schemes, "\n\n", "") // Remove extra newlines
	schemesSlice := strings.Split(schemes, "\n")

	for i, scheme := range schemesSlice {
		schemesSlice[i] = strings.TrimSpace(scheme)
	}

	return schemesSlice
}

// GetXcodeBuildSettings fetches build settings for a given path, scheme, and configuration.
//
// Parameters:
// - path (string): The project or workspace path.
// - schemeName (string): The scheme to use.
// - configuration (string): The build configuration (e.g., Debug, Release).
//
// Returns:
// - *XcodeBuildSettings: A struct containing the build settings.
// - error: If the settings cannot be retrieved or decoded.
func GetXcodeBuildSettings(path, schemeName, configuration string) (*XcodeBuildSettings, error) {
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

// getXcodeBuildSettings retrieves all build settings as a map from the `xcodebuild` output.
//
// Parameters:
// - path (string): The project or workspace path.
// - schemeName (string): The scheme to use.
// - configuration (string): The build configuration (optional).
//
// Returns:
// - *map[string]*string: A map of all build settings.
// - error: If the settings cannot be retrieved.
func getXcodeBuildSettings(path, schemeName, configuration string) (*map[string]*string, error) {
	var cmd *exec.Cmd

	if isXcodebuildInstalled() {
		if !strings.HasSuffix(path, ".xcworkspace") && !strings.HasSuffix(path, ".xcodeproj") {
			path = FindXcodeProjOrWorkspace(path)
		}

		var cmdArgs []string
		if strings.HasSuffix(path, ".xcworkspace") {
			cmdArgs = []string{"-workspace", path, "-scheme", schemeName, "-showBuildSettings"}
		} else if strings.HasSuffix(path, ".xcodeproj") {
			cmdArgs = []string{"-project", path, "-scheme", schemeName, "-showBuildSettings"}
		} else {
			return nil, fmt.Errorf("unable to locate .xcodeproj or .xcworkspace in the given path")
		}

		if configuration != "" {
			cmdArgs = append(cmdArgs, "-configuration", configuration)
		}

		cmd = exec.Command(utils.LocationOf(utils.XCODEBUILD), cmdArgs...)
	} else {
		return nil, fmt.Errorf("xcodebuild is not installed on this system")
	}

	output, err := cmd.Output()
	if err != nil {
		if strings.Contains(err.Error(), "exit status 65") {
			return nil, fmt.Errorf("scheme '%s' not found in location '%s'", schemeName, path)
		}
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

// IsPathAnXcodeProjectOrWorkspace checks if the given path is a .xcodeproj or .xcworkspace file.
//
// Parameters:
// - path (string): The path to check.
//
// Returns:
// - bool: True if the path is a valid Xcode project or workspace.
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

// GetDefaultProjectRoot determines the project root directory if none is provided.
//
// Parameters:
// - path (string): The current project path.
// - projectRoot (string): The explicitly specified project root (optional).
//
// Returns:
// - string: The resolved project root directory.
func GetDefaultProjectRoot(path, projectRoot string) string {
	if projectRoot == "" {
		if path == "" {
			currentDir, _ := os.Getwd()
			return currentDir
		}

		if utils.IsDir(path) {
			if strings.HasSuffix(path, ".xcodeproj") || strings.HasSuffix(path, ".xcworkspace") {
				return filepath.Dir(path)
			}
		}
		return path
	}
	return projectRoot
}

// isXcodebuildInstalled checks if the `xcodebuild` command is available on the system.
//
// Returns:
// - bool: True if `xcodebuild` is installed; false otherwise.
func isXcodebuildInstalled() bool {
	return utils.LocationOf(utils.XCODEBUILD) != ""
}

// FindXcodeProjOrWorkspace searches for a .xcodeproj or .xcworkspace file in the specified directory.
//
// Parameters:
// - path (string): The directory to search.
//
// Returns:
// - string: The path to the .xcodeproj or .xcworkspace file, preferring .xcworkspace if both are found.
func FindXcodeProjOrWorkspace(path string) string {
	var xcodeProjPath, xcodeWorkspacePath string

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
	}
	return xcodeProjPath
}
