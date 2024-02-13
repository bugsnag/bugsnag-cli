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
	DsymPaths   []string
}

// XcodeBuildSettings contains the relevant build settings required for uploading to bugsnag
type XcodeBuildSettings struct {
	ConfigurationBuildDir string `mapstructure:"CONFIGURATION_BUILD_DIR"`
	InfoPlistPath         string `mapstructure:"INFOPLIST_PATH"`
	BuiltProductsDir      string `mapstructure:"BUILT_PRODUCTS_DIR"`
	DsymName              string `mapstructure:"DWARF_DSYM_FILE_NAME"`
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
func GetXcodeBuildSettings(path, schemeName string) (*XcodeBuildSettings, error) {
	var buildSettings XcodeBuildSettings
	allBuildSettings, err := getXcodeBuildSettings(path, schemeName)
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
			cmd.Dir = path
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

// ProcessPathValues determines which values are set for each path value to be utilised downstream
func ProcessPathValues(path, dsymPath, projectRoot string) (*DsymUploadInfo, error) {

	// If dsymPath is set, then use it and don't set projectRoot
	if dsymPath != "" {
		foundDsymLocations, _ := findDsyms(dsymPath)
		return &DsymUploadInfo{"", foundDsymLocations}, nil
	}

	// If path is set
	if path != "" {

		// If path is also a directory
		if utils.IsDir(path) {

			// If projectRoot is set, use it for downstream
			if projectRoot != "" {
				log.Info("--project-root flag set, it's value takes precedence and will be used for upload")
				return &DsymUploadInfo{projectRoot, []string{""}}, nil
			}

			// If path is pointing to a .xcodeproj or .xcworkspace directory, set projectRoot to one directory up
			if strings.HasSuffix(path, ".xcodeproj") || strings.HasSuffix(path, ".xcworkspace") {
				// If path is pointing to a .xcodeproj or .xcworkspace directory, set projectRoot to one directory up
				return &DsymUploadInfo{filepath.Dir(path), []string{""}}, nil
			} else {

				// If path is a directory (not .xcodeproj or .xcworkspace), check for dSYMs within it
				foundDsymLocations, _ := findDsyms(path)

				if len(foundDsymLocations) != 0 {
					// If there are dSYMs found, then don't set projectRoot and set dsymPaths to the found dSYM locations
					return &DsymUploadInfo{"", foundDsymLocations}, nil
				} else {
					// If path is pointing to a directory and no dSYMs found within it, set projectRoot with path
					return &DsymUploadInfo{path, []string{""}}, nil
				}
			}

		} else {
			// If path is pointing to a file, we will assume it's pointing to a dSYM and use as-is
			return &DsymUploadInfo{"", []string{path}}, nil
		}

	}

	return nil, nil
}

//// ProcessPathValue determines the projectRoot from a given path
//func ProcessPathValue(path string, projectRoot string) (*DsymUploadInfo, error) {
//	if path == "" && projectRoot == "" {
//		currentDir, err := os.Getwd()
//		if err != nil {
//			return nil, err
//		}
//		return &DsymUploadInfo{currentDir, []string{""}}, nil
//	}
//
//	_, err := os.Stat(path)
//	if err != nil {
//		return nil, err
//	}
//
//	if utils.IsDir(path) {
//
//		if projectRoot != "" {
//			log.Info("--project-root flag set, it's value takes precedence and will be used for upload")
//			return &DsymUploadInfo{projectRoot, []string{""}}, nil
//		}
//
//		if strings.HasSuffix(path, ".xcodeproj") || strings.HasSuffix(path, ".xcworkspace") {
//			// If path is pointing to a .xcodeproj or .xcworkspace directory, set projectRoot to one directory up
//			return &DsymUploadInfo{filepath.Dir(path), []string{""}}, nil
//		} else {
//			dsymPaths, _ := findDsyms(path)
//			//fmt.Print(dsymPaths)
//
//			if len(dsymPaths) != 0 {
//				// Otherwise, don't set project root and set dsymPaths to the found dSYM locations
//				return &DsymUploadInfo{"", dsymPaths}, nil
//			} else {
//				// If path is pointing to a directory and no dSYMs found within it, set projectRoot to the path
//				return &DsymUploadInfo{path, []string{""}}, nil
//			}
//
//		}
//
//	} else {
//		// If path is pointing to a file, we will assume it's pointing to a dSYM and use as-is
//		return &DsymUploadInfo{"", []string{path}}, nil
//	}
//
//}

func findDsyms(root string) ([]string, error) {
	var dsyms []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.HasSuffix(info.Name(), ".dSYM") {
			dsyms = append(dsyms, filepath.Join(path, "Contents", "Resources", "DWARF"))
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return dsyms, nil
}

// isXcodebuildInstalled checks if xcodebuild is installed by checking if there is a path returned for it
func isXcodebuildInstalled() bool {
	return utils.LocationOf(utils.XCODEBUILD) != ""
}
