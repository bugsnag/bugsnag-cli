package unity

import (
	"fmt"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
	"path/filepath"
)

// GetAndroidLineMapping locates the LineNumberMappings.json file for Android builds.
//
// This function attempts to resolve the path to the IL2CPP line number mapping file,
// used for symbolication or debugging. The resolution follows this order:
//
//  1. If the input 'path' is non-empty, it is returned as-is.
//  2. It checks the default path: Library/Bee/artifacts/Android/il2cppOutput/cpp/Symbols/LineNumberMappings.json.
//  3. If not found, it searches under a backup folder (whose name ends with
//     "BackUpThisFolder_ButDontShipItWithYourGame" inside 'projectRoot') for the same file.
//
// Parameters:
//
//	path        - an optional explicit path to the mapping file. If provided, it is returned directly.
//	projectRoot - the root path of the Unity project, used to search for backup artifacts.
//
// Returns:
//
//	mappingPath - the resolved path to LineNumberMappings.json, or an empty string if not found.
//	error       - non-nil if there was an error during backup folder resolution.
func GetAndroidLineMapping(path string, buildDir string) (string, error) {
	if path != "" {
		return path, nil
	}

	// Check default artifacts path
	defaultPath := filepath.Join(buildDir, "Library", "Bee", "artifacts", "Android", "il2cppOutput", "cpp", "Symbols", "LineNumberMappings.json")
	if utils.FileExists(defaultPath) {
		return defaultPath, nil
	}

	backupDir, err := utils.FindFolderWithSuffix(buildDir, "BackUpThisFolder_ButDontShipItWithYourGame")
	if err != nil {
		return "", fmt.Errorf("unable to find backup folder: %s", err.Error())
	}

	backupPath := filepath.Join(backupDir, "il2cppOutput", "Symbols", "LineNumberMappings.json")
	if utils.FileExists(backupPath) {
		return backupPath, nil
	}

	return "", fmt.Errorf("Unable to fine line mapping file in your project: %s ", buildDir)
}
