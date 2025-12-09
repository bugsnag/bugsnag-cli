package android

import (
	"path/filepath"

	"github.com/bugsnag/bugsnag-cli/pkg/utils"
)

// FindAndroidManifest locates and returns the first existing AndroidManifest.xml
// for a given build directory and variant.
//
// This function constructs and checks multiple potential manifest locations that
// vary based on Android build outputs. It returns the first manifest file found
// on disk.
//
// Parameters:
//   - path: The Android merged_manifests directory.
//   - variant: The build variant (e.g., "debug", "release") whose manifest paths
//     should be searched.
//
// Returns:
//   - string: The resolved manifest path if found, otherwise an empty string.
func FindAndroidManifest(path string, variant string) string {

	primaryAppManifestPath := filepath.Join(path, variant, "AndroidManifest.xml")

	fallbackAppManifestPath := filepath.Join(path, variant, "process"+variant+"Manifest", "AndroidManifest.xml")

	paths := []string{primaryAppManifestPath, fallbackAppManifestPath}

	for _, path := range paths {
		if utils.FileExists(path) {
			return path
		}
	}

	return ""
}
