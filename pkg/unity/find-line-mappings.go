package unity

import (
	"fmt"
	"path/filepath"

	"github.com/bugsnag/bugsnag-cli/pkg/utils"
)

// GetAndroidLineMapping locates the LineNumberMappings.json file for Android builds.
//
// This function attempts to resolve the path to the IL2CPP line number mapping file,
// used for symbolication or debugging. The resolution follows this order:
//
//  1. It checks the default path: Library/Bee/artifacts/Android/il2cppOutput/cpp/Symbols/LineNumberMappings.json.
//  2. If not found, it searches under a backup folder (whose name ends with
//     "BackUpThisFolder_ButDontShipItWithYourGame" inside 'projectRoot') for the same file.
//
// Parameters:
//
//	projectRoot - the root path of the Unity project, used to search for backup artifacts.
//
// Returns:
//
//	mappingPath - the resolved path to LineNumberMappings.json, or an empty string if not found.
//	error       - non-nil if there was an error during backup folder resolution.
func GetAndroidLineMapping(projectRoot string) (string, error) {
	// Check the default artifacts path
	defaultPath := filepath.Join(projectRoot, "Library", "Bee", "artifacts", "Android", "il2cppOutput", "cpp", "Symbols", "LineNumberMappings.json")
	if utils.FileExists(defaultPath) {
		return defaultPath, nil
	}

	backupDir, err := utils.FindFolderWithSuffix(projectRoot, "BackUpThisFolder_ButDontShipItWithYourGame")
	if err != nil {
		return "", fmt.Errorf("unable to find backup folder: %s", err.Error())
	}

	backupPath := filepath.Join(backupDir, "il2cppOutput", "Symbols", "LineNumberMappings.json")
	if utils.FileExists(backupPath) {
		return backupPath, nil
	}

	return "", fmt.Errorf("Unable to fine line mapping file in your project: %s ", projectRoot)
}

// GetIosLineMapping locates the LineNumberMappings.json file for iOS builds.
//
// This function attempts to resolve the path to the IL2CPP line number mapping file,
// used for symbolication or debugging in iOS builds. The resolution follows this order:
//
//  1. It checks the default path: Library/Bee/artifacts/iOS/il2cppOutput/cpp/Symbols/LineNumberMappings.json.
//  2. If not found, it searches under a backup folder (ending with "_xcode" inside 'projectRoot') at:
//     Il2CppOutputProject/Source/il2cppOutput/Symbols/LineNumberMappings.json.
//
// Parameters:
//
//	projectRoot - the root path of the Unity project, used to search for backup artifacts.
//
// Returns:
//
//	mappingPath - the resolved path to LineNumberMappings.json.
//	error       - non-nil if the file cannot be found or the backup folder is missing.
func GetIosLineMapping(path string) (string, error) {
	var mappingPath string
	// Check the default artifacts path
	mappingPath = filepath.Join(path, "Library", "Bee", "artifacts", "iOS", "il2cppOutput", "cpp", "Symbols", "LineNumberMappings.json")
	if utils.FileExists(mappingPath) {
		return mappingPath, nil
	}

	// Try fallback: backup directory
	mappingPath = filepath.Join(path, "Il2CppOutputProject", "Source", "il2cppOutput", "Symbols", "LineNumberMappings.json")
	if utils.FileExists(mappingPath) {
		return mappingPath, nil
	}

	// Try fallback: temp directory
	mappingPath = filepath.Join(path, "Temp", "il2cppOutput", "il2cppOutput", "Symbols", "LineNumberMappings.json")
	if utils.FileExists(mappingPath) {
		return mappingPath, nil
	}

	// Unity 2021 Line Mapping Path Location
	mappingPath = filepath.Join(path, "Library", "Il2cppBuildCache", "iOS", "il2cppOutput", "Symbols", "LineNumberMappings.json")
	if utils.FileExists(mappingPath) {
		return mappingPath, nil
	}

	return "", fmt.Errorf("unable to find line mapping file in your project: %s", path)
}
